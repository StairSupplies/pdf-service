package handler

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/StairSupplies/pdf-service/internal/latex"
)

// Watermark handles POST /watermark.
//
// It reads the raw PDF from the request body, applies the requested watermark via
// a pdflatex pipeline, and returns the resulting PDF. All watermark parameters are
// supplied as request headers; the temp directory used during processing is always
// removed before the handler returns.
func Watermark(w http.ResponseWriter, r *http.Request) {
	text := r.Header.Get("X-Watermark-Text")
	if text == "" {
		log.Warn().Str("path", r.URL.Path).Msg("missing required header X-Watermark-Text")
		jsonError(w, "X-Watermark-Text header is required", http.StatusBadRequest)
		return
	}

	params := latex.Params{
		Text:     text,
		Color:    headerWithDefault(r, "X-Watermark-Color", "red"),
		Position: headerWithDefault(r, "X-Watermark-Position", "top-centre"),
		Bold:     r.Header.Get("X-Watermark-Bold") == "true",
	}

	var err error
	params.Opacity, err = parseFloat(headerWithDefault(r, "X-Watermark-Opacity", "0.5"))
	if err != nil {
		jsonError(w, "X-Watermark-Opacity must be a number", http.StatusBadRequest)
		return
	}

	sizeStr := headerWithDefault(r, "X-Watermark-Size", "60")
	params.Size, err = strconv.Atoi(sizeStr)
	if err != nil {
		jsonError(w, "X-Watermark-Size must be an integer", http.StatusBadRequest)
		return
	}

	params.Angle, err = parseFloat(headerWithDefault(r, "X-Watermark-Angle", "0"))
	if err != nil {
		jsonError(w, "X-Watermark-Angle must be a number", http.StatusBadRequest)
		return
	}

	// Read PDF body.
	pdfBytes, err := io.ReadAll(r.Body)
	if err != nil || len(pdfBytes) == 0 {
		jsonError(w, "request body must be a non-empty PDF", http.StatusBadRequest)
		return
	}

	// Create an isolated temp dir named /tmp/{uuid}; always cleaned up on return.
	tmpDir, err := uuidTempDir()
	if err != nil {
		log.Error().Err(err).Msg("failed to create temp dir")
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tmpDir)

	if err := os.WriteFile(filepath.Join(tmpDir, "input.pdf"), pdfBytes, 0644); err != nil {
		log.Error().Err(err).Msg("failed to write input.pdf")
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := latex.WriteJobTex(tmpDir, params); err != nil {
		log.Error().Err(err).Msg("failed to write job.tex")
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	pdflatexOut, err := latex.RunPdflatex(tmpDir)
	if err != nil {
		log.Error().Err(err).Msg("pdflatex failed")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(pdflatexOut) //nolint:errcheck
		return
	}

	outPDF, err := os.ReadFile(filepath.Join(tmpDir, "job.pdf"))
	if err != nil {
		log.Error().Err(err).Msg("failed to read job.pdf")
		jsonError(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Info().
		Str("text", params.Text).
		Str("color", params.Color).
		Float64("opacity", params.Opacity).
		Int("size", params.Size).
		Str("position", params.Position).
		Float64("angle", params.Angle).
		Int("output_bytes", len(outPDF)).
		Msg("watermark complete")

	w.Header().Set("Content-Type", "application/pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(outPDF) //nolint:errcheck
}

func jsonError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message}) //nolint:errcheck
}

func headerWithDefault(r *http.Request, key, def string) string {
	if v := r.Header.Get(key); v != "" {
		return v
	}
	return def
}

func parseFloat(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

// uuidTempDir creates and returns a temp directory at /tmp/{uuid v4}.
// Using a UUID ensures the path is unpredictable and avoids collisions.
func uuidTempDir() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("rand: %w", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant bits
	id := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	dir := filepath.Join(os.TempDir(), id)
	if err := os.Mkdir(dir, 0700); err != nil {
		return "", fmt.Errorf("mkdir: %w", err)
	}
	return dir, nil
}
