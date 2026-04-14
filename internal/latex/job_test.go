package latex_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/StairSupplies/pdf-service/internal/latex"
)

// readJobTex is a helper that writes a job.tex and returns its contents.
func readJobTex(t *testing.T, p latex.Params) string {
	t.Helper()
	dir := t.TempDir()
	if err := latex.WriteJobTex(dir, p); err != nil {
		t.Fatalf("WriteJobTex: %v", err)
	}
	b, err := os.ReadFile(filepath.Join(dir, "job.tex"))
	if err != nil {
		t.Fatalf("ReadFile job.tex: %v", err)
	}
	return string(b)
}

func TestWriteJobTex_RequiredPackages(t *testing.T) {
	tex := readJobTex(t, latex.Params{
		Text: "DRAFT", Color: "red", Opacity: 0.5, Size: 60,
	})
	for _, pkg := range []string{`\usepackage{pdfpages}`, `\usepackage{tikz}`} {
		if !strings.Contains(tex, pkg) {
			t.Errorf("expected %q in job.tex", pkg)
		}
	}
}

func TestWriteJobTex_InputPdfIncluded(t *testing.T) {
	tex := readJobTex(t, latex.Params{Text: "DRAFT", Color: "gray", Opacity: 0.3, Size: 40})
	if !strings.Contains(tex, "input.pdf") {
		t.Error("expected input.pdf referenced in job.tex")
	}
	if !strings.Contains(tex, `\includepdf`) {
		t.Error("expected \\includepdf in job.tex")
	}
}

func TestWriteJobTex_TextIsEscaped(t *testing.T) {
	tex := readJobTex(t, latex.Params{
		Text: `100% off $10`, Color: "red", Opacity: 0.5, Size: 60,
	})
	// Raw special chars must not appear; escaped forms must.
	if strings.Contains(tex, `100% off $10`) {
		t.Error("watermark text was not escaped in job.tex")
	}
	if !strings.Contains(tex, `100\% off \$10`) {
		t.Error("expected escaped watermark text in job.tex")
	}
}

func TestWriteJobTex_ColorAndOpacity(t *testing.T) {
	tex := readJobTex(t, latex.Params{
		Text: "X", Color: "blue", Opacity: 0.75, Size: 48,
	})
	if !strings.Contains(tex, "text=blue") {
		t.Error("expected text=blue in job.tex")
	}
	if !strings.Contains(tex, "opacity=0.75") {
		t.Error("expected opacity=0.75 in job.tex")
	}
}

func TestWriteJobTex_UnsafeColorFallsBackToBlack(t *testing.T) {
	tex := readJobTex(t, latex.Params{
		Text: "X", Color: `red; \evil`, Opacity: 0.5, Size: 60,
	})
	if strings.Contains(tex, `\evil`) {
		t.Error("unsafe color was not sanitised")
	}
	if !strings.Contains(tex, "text=black") {
		t.Error("expected fallback text=black for unsafe input")
	}
}

var positionTests = []struct {
	position   string
	wantAnchor string
	wantVShift string
}{
	{"top-centre", "current page.north", "-2cm"},
	{"top-center", "current page.north", "-2cm"},
	{"top-left", "current page.north west", "-2cm"},
	{"top-right", "current page.north east", "-2cm"},
	{"bottom-centre", "current page.south", "2cm"},
	{"bottom-center", "current page.south", "2cm"},
	{"bottom-left", "current page.south west", "2cm"},
	{"bottom-right", "current page.south east", "2cm"},
	{"centre", "current page.center", "0pt"},
	{"center", "current page.center", "0pt"},
	{"unknown", "current page.center", "0pt"},
}

func TestWriteJobTex_Positions(t *testing.T) {
	for _, tc := range positionTests {
		t.Run(tc.position, func(t *testing.T) {
			tex := readJobTex(t, latex.Params{
				Text: "X", Color: "red", Opacity: 0.5, Size: 60, Position: tc.position,
			})
			if !strings.Contains(tex, tc.wantAnchor) {
				t.Errorf("position %q: expected anchor %q in job.tex\ngot:\n%s", tc.position, tc.wantAnchor, tex)
			}
			if !strings.Contains(tex, tc.wantVShift) {
				t.Errorf("position %q: expected vshift %q in job.tex\ngot:\n%s", tc.position, tc.wantVShift, tex)
			}
		})
	}
}

func TestWriteJobTex_AngleAndSize(t *testing.T) {
	tex := readJobTex(t, latex.Params{
		Text: "X", Color: "red", Opacity: 0.5, Size: 72, Angle: 45,
	})
	if !strings.Contains(tex, "rotate=45") {
		t.Error("expected rotate=45 in job.tex")
	}
	if !strings.Contains(tex, "72") {
		t.Error("expected font size 72 in job.tex")
	}
}

func TestWriteJobTex_Bold(t *testing.T) {
	t.Run("bold enabled", func(t *testing.T) {
		tex := readJobTex(t, latex.Params{Text: "X", Color: "red", Opacity: 0.5, Size: 60, Bold: true})
		if !strings.Contains(tex, `\bfseries`) {
			t.Error("expected \\bfseries in job.tex when Bold=true")
		}
	})
	t.Run("bold disabled", func(t *testing.T) {
		tex := readJobTex(t, latex.Params{Text: "X", Color: "red", Opacity: 0.5, Size: 60, Bold: false})
		if strings.Contains(tex, `\bfseries`) {
			t.Error("unexpected \\bfseries in job.tex when Bold=false")
		}
	})
}

func TestWriteJobTex_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	if err := latex.WriteJobTex(dir, latex.Params{Text: "X", Color: "red", Opacity: 0.5, Size: 60}); err != nil {
		t.Fatalf("WriteJobTex returned error: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "job.tex")); err != nil {
		t.Errorf("job.tex not created: %v", err)
	}
}
