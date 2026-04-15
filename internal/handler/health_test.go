package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/StairSupplies/pdf-service/internal/handler"
)

func TestHealth_OK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.Health(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want status 200, got %d", rr.Code)
	}

	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Errorf("want Content-Type application/json, got %q", ct)
	}

	if !strings.Contains(rr.Body.String(), "ok") {
		t.Errorf("want body to contain %q, got %q", "ok", rr.Body.String())
	}
}
