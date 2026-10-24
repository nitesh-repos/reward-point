package server

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reward-points/internal/database"
	"testing"
)

func newTestServer() http.Handler {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return New(Config{}, logger, database.NewLocalStorage()).httpServer.Handler
}

func TestGetMember_NotFound(t *testing.T) {
	h := newTestServer()
	req := httptest.NewRequest(http.MethodGet, "/members/42", nil)
	rec := httptest.NewRecorder()

	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}
