package latex

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"text/template"
)

// Params holds watermark rendering parameters for a pdflatex job.
type Params struct {
	Text     string  // raw watermark text (escaped internally before use in TeX)
	Color    string  // xcolor color name (e.g. "red", "gray")
	Opacity  float64 // 0.0–1.0
	Size     int     // font size in pt
	Position string  // "top-centre", "centre", "bottom-centre", etc.
	Angle    float64 // rotation in degrees
	Bold     bool    // render text in bold weight
}

// safeColorRe matches only characters that are safe inside a LaTeX xcolor argument,
// preventing TeX injection via the color parameter.
var safeColorRe = regexp.MustCompile(`^[a-zA-Z0-9!.,{} %-]+$`)

// texSrc is the pdflatex template. Delimiters are << >> to avoid conflicts with LaTeX braces.
// TikZ overlay is used instead of the background package so the watermark renders on top of
// the included PDF content (the background package renders behind page content, making it
// invisible on PDFs with a solid white background).
const texSrc = `\documentclass{article}
\usepackage{lmodern}
\usepackage{pdfpages}
\usepackage{tikz}
\begin{document}
\includepdf[pages=-,fitpaper,pagecommand={%
  \begin{tikzpicture}[remember picture,overlay]
    \node[rotate=<<.Angle>>,opacity=<<.Opacity>>,text=<<.Color>>,
          font=<<.BoldCmd>>\fontsize{<<.Size>>}{<<.Size>>}\selectfont,
          yshift=<<.VShift>>]
      at (<<.Position>>) {<<.Text>>};
  \end{tikzpicture}%
}]{input.pdf}
\end{document}
`

var tmpl = template.Must(template.New("job").Delims("<<", ">>").Parse(texSrc))

type jobData struct {
	Text     string
	Color    string
	Opacity  float64
	Size     int
	Angle    float64
	Position string
	VShift   string
	BoldCmd  string // "\bfseries" or ""
}

// positionAndShift maps the human-readable position string to a TikZ anchor and
// a vshift value that pulls the text away from the page edge.
func positionAndShift(pos string) (tikzPos, vshift string) {
	switch pos {
	case "top-left":
		return "current page.north west", "-2cm"
	case "top-right":
		return "current page.north east", "-2cm"
	case "top-centre", "top-center":
		return "current page.north", "-2cm"
	case "bottom-left":
		return "current page.south west", "2cm"
	case "bottom-right":
		return "current page.south east", "2cm"
	case "bottom-centre", "bottom-center":
		return "current page.south", "2cm"
	default: // "centre", "center", unrecognised
		return "current page.center", "0pt"
	}
}

// WriteJobTex writes job.tex into dir using the provided Params.
// The color is sanitised against a safe allowlist; an unsafe value falls back to "black".
func WriteJobTex(dir string, p Params) error {
	color := p.Color
	if !safeColorRe.MatchString(color) {
		color = "black"
	}

	pos, vshift := positionAndShift(p.Position)

	boldCmd := ""
	if p.Bold {
		boldCmd = `\bfseries`
	}

	data := jobData{
		Text:     Escape(p.Text),
		Color:    color,
		Opacity:  p.Opacity,
		Size:     p.Size,
		Angle:    p.Angle,
		Position: pos,
		VShift:   vshift,
		BoldCmd:  boldCmd,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("template: %w", err)
	}
	return os.WriteFile(filepath.Join(dir, "job.tex"), buf.Bytes(), 0644)
}

// RunPdflatex runs pdflatex on job.tex inside dir (also used as cwd and output dir).
// Two passes are required: the first pass writes TikZ overlay coordinates to the .aux
// file; the second pass reads them to correctly position the watermark on the page.
// On failure it returns the combined stdout+stderr together with the error so the
// caller can relay it to the client.
func RunPdflatex(dir string) ([]byte, error) {
	for range 2 {
		cmd := exec.Command(
			"pdflatex",
			"-interaction=nonstopmode",
			"-output-directory", dir,
			"job.tex",
		)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			return out, fmt.Errorf("pdflatex: %w", err)
		}
	}
	return nil, nil
}
