package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"strings"
	"testing"

	"github.com/StairSupplies/pdf-service/internal/handler"
)

// newWatermarkRequest builds a POST /watermark request with the given body and headers.
func newWatermarkRequest(body []byte, headers map[string]string) *http.Request {
	req := httptest.NewRequest(http.MethodPost, "/watermark", bytes.NewReader(body))
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	return req
}

func TestWatermark_MissingText(t *testing.T) {
	req := newWatermarkRequest([]byte("%PDF-1.4 fake"), nil)
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rr.Code)
	}
}

func TestWatermark_InvalidOpacity(t *testing.T) {
	req := newWatermarkRequest([]byte("%PDF-1.4 fake"), map[string]string{
		"X-Watermark-Text":    "DRAFT",
		"X-Watermark-Opacity": "notanumber",
	})
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rr.Code)
	}
}

func TestWatermark_InvalidSize(t *testing.T) {
	req := newWatermarkRequest([]byte("%PDF-1.4 fake"), map[string]string{
		"X-Watermark-Text": "DRAFT",
		"X-Watermark-Size": "big",
	})
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rr.Code)
	}
}

func TestWatermark_InvalidAngle(t *testing.T) {
	req := newWatermarkRequest([]byte("%PDF-1.4 fake"), map[string]string{
		"X-Watermark-Text":  "DRAFT",
		"X-Watermark-Angle": "sideways",
	})
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rr.Code)
	}
}

func TestWatermark_EmptyBody(t *testing.T) {
	req := newWatermarkRequest(nil, map[string]string{
		"X-Watermark-Text": "DRAFT",
	})
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Errorf("want 400, got %d", rr.Code)
	}
}

func TestWatermark_BodyTooLarge(t *testing.T) {
	// 121 MB — just over the 120 MB limit.
	oversized := bytes.Repeat([]byte("x"), 121<<20)
	req := newWatermarkRequest(oversized, map[string]string{
		"X-Watermark-Text": "DRAFT",
	})
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)
	if rr.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("want 413, got %d", rr.Code)
	}
}

// TestWatermark_ContentType verifies that a valid request returns an application/pdf
// response. This test requires pdflatex to be installed and is skipped when it is not.
func TestWatermark_ContentType(t *testing.T) {
	if _, err := exec.LookPath("pdflatex"); err != nil {
		t.Skip("pdflatex not installed; skipping integration test")
	}

	// Minimal valid single-page PDF (version 1.4, hand-crafted).
	const minimalPDF = "%PDF-1.4\n" +
		"1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n" +
		"2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n" +
		"3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 612 792] >>\nendobj\n" +
		"xref\n0 4\n0000000000 65535 f \n0000000009 00000 n \n0000000058 00000 n \n" +
		"0000000115 00000 n \n" +
		"trailer\n<< /Size 4 /Root 1 0 R >>\nstartxref\n190\n%%EOF\n"

	req := newWatermarkRequest([]byte(minimalPDF), map[string]string{
		"X-Watermark-Text": "REMAKE: 653374-02",
	})
	rr := httptest.NewRecorder()
	handler.Watermark(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("want 200, got %d — body: %s", rr.Code, rr.Body.String())
	}

	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/pdf") {
		t.Errorf("want Content-Type application/pdf, got %q", ct)
	}
}
