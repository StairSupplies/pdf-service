package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

	"github.com/StairSupplies/pdf-service/internal/latex"
)

// WatermarkParams holds validated parameters parsed from request headers.
type WatermarkParams struct {
	Text     string
	Color    string
	Opacity  float64
	Size     int
	Position string
	Angle    float64
}

// Watermark handles POST /watermark — parses headers, validates required fields,
// and (in Piece 1) returns 202 Accepted with the LaTeX-escaped watermark text.
func Watermark(w http.ResponseWriter, r *http.Request) {
	text := r.Header.Get("X-Watermark-Text")
	if text == "" {
		log.Warn().Str("path", r.URL.Path).Msg("missing required header X-Watermark-Text")
		http.Error(w, `{"error":"X-Watermark-Text header is required"}`, http.StatusBadRequest)
		return
	}

	params := WatermarkParams{
		Text:     text,
		Color:    headerWithDefault(r, "X-Watermark-Color", "red"),
		Position: headerWithDefault(r, "X-Watermark-Position", "top-centre"),
	}

	var err error
	params.Opacity, err = parseFloat(headerWithDefault(r, "X-Watermark-Opacity", "0.5"))
	if err != nil {
		http.Error(w, `{"error":"X-Watermark-Opacity must be a number"}`, http.StatusBadRequest)
		return
	}

	sizeStr := headerWithDefault(r, "X-Watermark-Size", "60")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		http.Error(w, `{"error":"X-Watermark-Size must be an integer"}`, http.StatusBadRequest)
		return
	}
	params.Size = size

	params.Angle, err = parseFloat(headerWithDefault(r, "X-Watermark-Angle", "0"))
	if err != nil {
		http.Error(w, `{"error":"X-Watermark-Angle must be a number"}`, http.StatusBadRequest)
		return
	}

	escapedText := latex.Escape(params.Text)

	log.Info().
		Str("text", params.Text).
		Str("color", params.Color).
		Float64("opacity", params.Opacity).
		Int("size", params.Size).
		Str("position", params.Position).
		Float64("angle", params.Angle).
		Msg("watermark request accepted")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{ //nolint:errcheck
		"status": "accepted",
		"text":   escapedText,
	})
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
