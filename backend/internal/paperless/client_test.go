package paperless

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// fakePaperless is a minimal paperless-ngx API double.
type fakePaperless struct {
	mux          *http.ServeMux
	server       *httptest.Server
	customFields []customField
	documents    map[int]*fakeDocument
	lastAuth     string
	lastAccept   string
}

type fakeDocument struct {
	ID           int                `json:"id"`
	Title        string             `json:"title"`
	Created      string             `json:"created"`
	CustomFields []customFieldValue `json:"custom_fields"`
}

func newFakePaperless(t *testing.T) *fakePaperless {
	t.Helper()
	f := &fakePaperless{
		mux:       http.NewServeMux(),
		documents: map[int]*fakeDocument{},
	}
	record := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			f.lastAuth = r.Header.Get("Authorization")
			f.lastAccept = r.Header.Get("Accept")
			next(w, r)
		}
	}
	f.mux.HandleFunc("GET /api/documents/", record(f.listDocuments))
	f.mux.HandleFunc("GET /api/documents/{id}/", record(f.getDocument))
	f.mux.HandleFunc("PATCH /api/documents/{id}/", record(f.patchDocument))
	f.mux.HandleFunc("GET /api/documents/{id}/thumb/", record(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write([]byte("thumb-bytes"))
	}))
	f.mux.HandleFunc("GET /api/custom_fields/", record(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(customFieldPage{Count: len(f.customFields), Results: f.customFields})
	}))
	f.mux.HandleFunc("POST /api/custom_fields/", record(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		field := customField{ID: len(f.customFields) + 1, Name: body["name"], DataType: body["data_type"]}
		f.customFields = append(f.customFields, field)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(field)
	}))
	f.server = httptest.NewServer(f.mux)
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakePaperless) listDocuments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	results := []map[string]any{}
	for _, doc := range f.documents {
		entry := map[string]any{"id": doc.ID, "title": doc.Title, "created": doc.Created}
		if query != "" {
			entry["__search_hit__"] = map[string]any{"highlights": "<span>match</span> context"}
		}
		results = append(results, entry)
	}
	json.NewEncoder(w).Encode(map[string]any{"count": len(results), "results": results})
}

func (f *fakePaperless) getDocument(w http.ResponseWriter, r *http.Request) {
	doc, ok := f.documents[atoiOr(r.PathValue("id"))]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(doc)
}

func (f *fakePaperless) patchDocument(w http.ResponseWriter, r *http.Request) {
	doc, ok := f.documents[atoiOr(r.PathValue("id"))]
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	var body documentFields
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "bad body", http.StatusBadRequest)
		return
	}
	doc.CustomFields = body.CustomFields
	json.NewEncoder(w).Encode(doc)
}

func atoiOr(s string) int {
	var n int
	fmt.Sscanf(s, "%d", &n)
	return n
}

func (f *fakePaperless) client() *Client {
	return NewClient(f.server.URL, "test-token")
}

func TestClientSendsAuthAndVersionHeaders(t *testing.T) {
	f := newFakePaperless(t)
	if err := f.client().Ping(context.Background()); err != nil {
		t.Fatalf("Ping: %v", err)
	}
	if f.lastAuth != "Token test-token" {
		t.Errorf("Authorization = %q, want %q", f.lastAuth, "Token test-token")
	}
	if f.lastAccept != "application/json; version="+apiVersion {
		t.Errorf("Accept = %q", f.lastAccept)
	}
}

func TestSearchDocuments(t *testing.T) {
	f := newFakePaperless(t)
	f.documents[1] = &fakeDocument{ID: 1, Title: "Invoice 2026", Created: "2026-01-01"}

	page, err := f.client().SearchDocuments(context.Background(), "invoice", 1)
	if err != nil {
		t.Fatalf("SearchDocuments: %v", err)
	}
	if page.Count != 1 || len(page.Results) != 1 {
		t.Fatalf("page = %+v", page)
	}
	res := page.Results[0]
	if res.Title != "Invoice 2026" {
		t.Errorf("Title = %q", res.Title)
	}
	if res.Snippet != "match context" {
		t.Errorf("Snippet = %q, want HTML stripped %q", res.Snippet, "match context")
	}
}

func TestSearchDocuments_EmptyQueryListsRecent(t *testing.T) {
	f := newFakePaperless(t)
	f.documents[1] = &fakeDocument{ID: 1, Title: "Doc", Created: "2026-01-01"}

	page, err := f.client().SearchDocuments(context.Background(), "", 1)
	if err != nil {
		t.Fatalf("SearchDocuments: %v", err)
	}
	if len(page.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(page.Results))
	}
	if page.Results[0].Snippet != "" {
		t.Errorf("Snippet = %q, want empty for non-search listing", page.Results[0].Snippet)
	}
}

func TestEnsureCustomField_CreatesThenFinds(t *testing.T) {
	f := newFakePaperless(t)
	c := f.client()

	id, err := c.EnsureCustomField(context.Background())
	if err != nil {
		t.Fatalf("EnsureCustomField: %v", err)
	}
	if id != 1 {
		t.Fatalf("id = %d, want 1", id)
	}
	if len(f.customFields) != 1 || f.customFields[0].DataType != "url" {
		t.Fatalf("customFields = %+v", f.customFields)
	}

	id2, err := c.EnsureCustomField(context.Background())
	if err != nil {
		t.Fatalf("EnsureCustomField (second): %v", err)
	}
	if id2 != id {
		t.Errorf("second call id = %d, want %d (no duplicate)", id2, id)
	}
	if len(f.customFields) != 1 {
		t.Errorf("expected 1 custom field, got %d", len(f.customFields))
	}
}

func TestSetCustomField_PreservesOtherFields(t *testing.T) {
	f := newFakePaperless(t)
	f.documents[5] = &fakeDocument{ID: 5, Title: "Doc", CustomFields: []customFieldValue{{Field: 9, Value: "keep"}}}

	if err := f.client().SetCustomField(context.Background(), 5, 1, "https://kontor.example/a"); err != nil {
		t.Fatalf("SetCustomField: %v", err)
	}
	fields := f.documents[5].CustomFields
	if len(fields) != 2 {
		t.Fatalf("fields = %+v", fields)
	}
	if err := f.client().SetCustomField(context.Background(), 5, 1, "https://kontor.example/b"); err != nil {
		t.Fatalf("SetCustomField (update): %v", err)
	}
	fields = f.documents[5].CustomFields
	if len(fields) != 2 {
		t.Fatalf("expected value update not append, fields = %+v", fields)
	}
	for _, fv := range fields {
		if fv.Field == 1 && fv.Value != "https://kontor.example/b" {
			t.Errorf("field 1 value = %v", fv.Value)
		}
	}
}

func TestClearCustomFieldIfMatches(t *testing.T) {
	f := newFakePaperless(t)
	f.documents[5] = &fakeDocument{ID: 5, CustomFields: []customFieldValue{
		{Field: 1, Value: "https://kontor.example/other"},
		{Field: 9, Value: "keep"},
	}}

	// Non-matching value: field stays.
	if err := f.client().ClearCustomFieldIfMatches(context.Background(), 5, 1, "https://kontor.example/mine"); err != nil {
		t.Fatalf("ClearCustomFieldIfMatches: %v", err)
	}
	if len(f.documents[5].CustomFields) != 2 {
		t.Fatalf("field should not be cleared on mismatch: %+v", f.documents[5].CustomFields)
	}

	// Matching value: field removed, others kept.
	if err := f.client().ClearCustomFieldIfMatches(context.Background(), 5, 1, "https://kontor.example/other"); err != nil {
		t.Fatalf("ClearCustomFieldIfMatches: %v", err)
	}
	fields := f.documents[5].CustomFields
	if len(fields) != 1 || fields[0].Field != 9 {
		t.Fatalf("fields = %+v", fields)
	}
}

func TestClientUpstreamError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "invalid token", http.StatusUnauthorized)
	}))
	defer server.Close()

	err := NewClient(server.URL, "bad").Ping(context.Background())
	if err == nil {
		t.Fatal("expected error for 401 response")
	}
}
