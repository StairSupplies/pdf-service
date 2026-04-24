package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/StairSupplies/pdf-service/internal/middleware"
)

func okHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuth_KeySet_CorrectToken(t *testing.T) {
	h := middleware.Auth("secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/watermark", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rr.Code)
	}
}

func TestAuth_KeySet_MissingHeader(t *testing.T) {
	h := middleware.Auth("secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/watermark", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rr.Code)
	}
}

func TestAuth_KeySet_WrongToken(t *testing.T) {
	h := middleware.Auth("secret")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/watermark", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("want 401, got %d", rr.Code)
	}
}

func TestAuth_KeyUnset_PassThrough(t *testing.T) {
	h := middleware.Auth("")(okHandler())
	req := httptest.NewRequest(http.MethodPost, "/watermark", nil)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("want 200, got %d", rr.Code)
	}
}
