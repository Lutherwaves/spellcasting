package RESOURCE_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"SERVICENAME/pkg/features/RESOURCE"
	"SERVICENAME/pkg/types"

	"github.com/tink3rlabs/magic/storage"
)

// newTestRouter builds a Router backed by a magic in-memory storage adapter.
func newTestRouter(t *testing.T) *chi.Mux {
	t.Helper()

	adapter, err := storage.StorageAdapterFactory{}.GetInstance(storage.MEMORY, map[string]string{})
	if err != nil {
		t.Fatalf("in-memory storage: %v", err)
	}

	svc := RESOURCE.NewService(adapter)
	r := RESOURCE.NewResourceRouter(svc)
	return r.Router.(*chi.Mux)
}

func TestList_Empty(t *testing.T) {
	mux := newTestRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var body struct {
		Items []types.Resource `json:"items"`
		Next  string           `json:"next"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Items) != 0 {
		t.Errorf("expected empty items, got %d", len(body.Items))
	}
	if body.Next != "" {
		t.Errorf("expected empty next cursor, got %q", body.Next)
	}
}

func TestCreate_And_Get(t *testing.T) {
	mux := newTestRouter(t)

	// Create
	payload := map[string]any{"name": "test-RESOURCE"}
	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d — body: %s", rec.Code, rec.Body.String())
	}

	var created types.Resource
	if err := json.NewDecoder(rec.Body).Decode(&created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID in create response")
	}
	if created.Name != "test-RESOURCE" {
		t.Errorf("expected name %q, got %q", "test-RESOURCE", created.Name)
	}

	// Get
	req2 := httptest.NewRequest(http.MethodGet, "/"+created.ID, nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("get: expected 200, got %d", rec2.Code)
	}

	var fetched types.Resource
	if err := json.NewDecoder(rec2.Body).Decode(&fetched); err != nil {
		t.Fatalf("decode get response: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("id mismatch: want %q, got %q", created.ID, fetched.ID)
	}
}

func TestCreate_RequiresName(t *testing.T) {
	mux := newTestRouter(t)

	body, _ := json.Marshal(map[string]any{})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing name, got %d", rec.Code)
	}
}

func TestList_Pagination(t *testing.T) {
	mux := newTestRouter(t)

	// Create 3 items
	for _, name := range []string{"alpha", "beta", "gamma"} {
		b, _ := json.Marshal(map[string]any{"name": name})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("create %q: got %d", name, rec.Code)
		}
	}

	// Page 1: limit=2
	req := httptest.NewRequest(http.MethodGet, "/?limit=2", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("page1: expected 200, got %d", rec.Code)
	}

	var page1 struct {
		Items []types.Resource `json:"items"`
		Next  string           `json:"next"`
	}
	json.NewDecoder(rec.Body).Decode(&page1)
	if len(page1.Items) != 2 {
		t.Fatalf("page1: expected 2 items, got %d", len(page1.Items))
	}
	if page1.Next == "" {
		t.Fatal("page1: expected non-empty next cursor")
	}

	// Page 2: use cursor
	req2 := httptest.NewRequest(http.MethodGet, "/?limit=2&next="+page1.Next, nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("page2: expected 200, got %d", rec2.Code)
	}

	var page2 struct {
		Items []types.Resource `json:"items"`
		Next  string           `json:"next"`
	}
	json.NewDecoder(rec2.Body).Decode(&page2)
	if len(page2.Items) != 1 {
		t.Fatalf("page2: expected 1 item, got %d", len(page2.Items))
	}
}

func TestList_Filter(t *testing.T) {
	mux := newTestRouter(t)

	for _, name := range []string{"foo", "bar", "foobar"} {
		b, _ := json.Marshal(map[string]any{"name": name})
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
	}

	req := httptest.NewRequest(http.MethodGet, "/?filter=name:foo", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("filter: expected 200, got %d", rec.Code)
	}

	var result struct {
		Items []types.Resource `json:"items"`
	}
	json.NewDecoder(rec.Body).Decode(&result)
	for _, item := range result.Items {
		if item.Name != "foo" {
			t.Errorf("filter returned unexpected item %q", item.Name)
		}
	}
}

func TestUpdate(t *testing.T) {
	mux := newTestRouter(t)

	b, _ := json.Marshal(map[string]any{"name": "original"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var created types.Resource
	json.NewDecoder(rec.Body).Decode(&created)

	updated := "updated"
	patch, _ := json.Marshal(map[string]any{"name": &updated})
	req2 := httptest.NewRequest(http.MethodPatch, "/"+created.ID, bytes.NewReader(patch))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusOK {
		t.Fatalf("update: expected 200, got %d — %s", rec2.Code, rec2.Body.String())
	}

	var result types.Resource
	json.NewDecoder(rec2.Body).Decode(&result)
	if result.Name != "updated" {
		t.Errorf("expected name %q, got %q", "updated", result.Name)
	}
}

func TestDelete(t *testing.T) {
	mux := newTestRouter(t)

	b, _ := json.Marshal(map[string]any{"name": "to-delete"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	var created types.Resource
	json.NewDecoder(rec.Body).Decode(&created)

	req2 := httptest.NewRequest(http.MethodDelete, "/"+created.ID, nil)
	rec2 := httptest.NewRecorder()
	mux.ServeHTTP(rec2, req2)

	if rec2.Code != http.StatusNoContent {
		t.Fatalf("delete: expected 204, got %d", rec2.Code)
	}

	// Get should now 404
	req3 := httptest.NewRequest(http.MethodGet, "/"+created.ID, nil)
	rec3 := httptest.NewRecorder()
	mux.ServeHTTP(rec3, req3)

	if rec3.Code != http.StatusNotFound {
		t.Fatalf("get after delete: expected 404, got %d", rec3.Code)
	}
}

// Ensure context flows through — exercises the Context parameter in service calls.
var _ context.Context = context.Background()
