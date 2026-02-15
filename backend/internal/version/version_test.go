package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	Version = "v1.2.3"
	Commit = "abc1234"
	BuildDate = "2026-01-15"
	t.Cleanup(func() {
		Version = "dev"
		Commit = ""
		BuildDate = ""
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %s", ct)
	}

	var info Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if info.Version != "v1.2.3" {
		t.Errorf("version = %q, want %q", info.Version, "v1.2.3")
	}
	if info.Commit != "abc1234" {
		t.Errorf("commit = %q, want %q", info.Commit, "abc1234")
	}
	if info.BuildDate != "2026-01-15" {
		t.Errorf("buildDate = %q, want %q", info.BuildDate, "2026-01-15")
	}
}

func TestHandlerDefaults(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	Handler(rec, req)

	var info Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if info.Version != "dev" {
		t.Errorf("version = %q, want %q", info.Version, "dev")
	}
	if info.Commit != "" {
		t.Errorf("commit = %q, want empty", info.Commit)
	}
	if info.BuildDate != "" {
		t.Errorf("buildDate = %q, want empty", info.BuildDate)
	}
}
