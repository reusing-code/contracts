package version

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/version", nil)

	Handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	ct := rec.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %q, want %q", ct, "application/json")
	}
}

func TestHandler_ReturnsVersionInfo(t *testing.T) {
	// Set version variables
	Version = "1.2.3"
	Commit = "abc123"
	BuildDate = "2026-01-15"

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/version", nil)

	Handler(rec, req)

	var info Info
	if err := json.NewDecoder(rec.Body).Decode(&info); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if info.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", info.Version, "1.2.3")
	}
	if info.Commit != "abc123" {
		t.Errorf("Commit = %q, want %q", info.Commit, "abc123")
	}
	if info.BuildDate != "2026-01-15" {
		t.Errorf("BuildDate = %q, want %q", info.BuildDate, "2026-01-15")
	}
}

func TestHandler_OmitsEmptyFields(t *testing.T) {
	// Set only version
	Version = "dev"
	Commit = ""
	BuildDate = ""

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/version", nil)

	Handler(rec, req)

	var raw map[string]any
	if err := json.NewDecoder(rec.Body).Decode(&raw); err != nil {
		t.Fatalf("decoding response: %v", err)
	}

	if _, ok := raw["commit"]; ok {
		t.Error("response should not contain empty commit field")
	}
	if _, ok := raw["buildDate"]; ok {
		t.Error("response should not contain empty buildDate field")
	}
	if raw["version"] != "dev" {
		t.Errorf("version = %v, want %q", raw["version"], "dev")
	}
}

func TestGet_ReturnsCurrentVersionInfo(t *testing.T) {
	Version = "test-version"
	Commit = "test-commit"
	BuildDate = "test-date"

	info := Get()

	if info.Version != "test-version" {
		t.Errorf("Version = %q, want %q", info.Version, "test-version")
	}
	if info.Commit != "test-commit" {
		t.Errorf("Commit = %q, want %q", info.Commit, "test-commit")
	}
	if info.BuildDate != "test-date" {
		t.Errorf("BuildDate = %q, want %q", info.BuildDate, "test-date")
	}
}
