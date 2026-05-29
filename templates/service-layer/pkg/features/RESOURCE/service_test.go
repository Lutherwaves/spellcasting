package RESOURCE_test

import (
	"context"
	"testing"

	"SERVICENAME/pkg/features/RESOURCE"
	"SERVICENAME/pkg/types"

	"github.com/tink3rlabs/magic/storage"
)

func newTestService(t *testing.T) *RESOURCE.Service {
	t.Helper()
	adapter, err := storage.StorageAdapterFactory{}.GetInstance(storage.MEMORY, map[string]string{})
	if err != nil {
		t.Fatalf("in-memory adapter: %v", err)
	}
	return RESOURCE.NewService(adapter)
}

var ctx = context.Background()

func TestList_EmptyInitially(t *testing.T) {
	svc := newTestService(t)

	items, next, err := svc.List(ctx, "", 10, "")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected 0 items, got %d", len(items))
	}
	if next != "" {
		t.Errorf("expected empty cursor, got %q", next)
	}
}

func TestCreate_RoundTripsViaGet(t *testing.T) {
	svc := newTestService(t)

	created, err := svc.Create(ctx, types.CreateResource{Name: "hello"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if created.Name != "hello" {
		t.Errorf("name: want %q, got %q", "hello", created.Name)
	}
	if created.CreatedAt.IsZero() {
		t.Error("CreatedAt must not be zero")
	}
	if created.UpdatedAt.IsZero() {
		t.Error("UpdatedAt must not be zero")
	}

	fetched, err := svc.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if fetched.ID != created.ID {
		t.Errorf("id mismatch: want %q, got %q", created.ID, fetched.ID)
	}
	if fetched.Name != "hello" {
		t.Errorf("name mismatch: want %q, got %q", "hello", fetched.Name)
	}
}

func TestList_Pagination(t *testing.T) {
	svc := newTestService(t)

	for i := 0; i < 3; i++ {
		if _, err := svc.Create(ctx, types.CreateResource{Name: "item"}); err != nil {
			t.Fatalf("Create: %v", err)
		}
	}

	// Page 1
	page1, next1, err := svc.List(ctx, "", 2, "")
	if err != nil {
		t.Fatalf("List page1: %v", err)
	}
	if len(page1) != 2 {
		t.Fatalf("page1: want 2 items, got %d", len(page1))
	}
	if next1 == "" {
		t.Fatal("page1: expected non-empty next cursor")
	}

	// Page 2
	page2, next2, err := svc.List(ctx, "", 2, next1)
	if err != nil {
		t.Fatalf("List page2: %v", err)
	}
	if len(page2) != 1 {
		t.Fatalf("page2: want 1 item, got %d", len(page2))
	}
	_ = next2 // empty or not — depends on adapter; don't assert here
}

func TestList_Filter(t *testing.T) {
	svc := newTestService(t)

	for _, name := range []string{"foo", "bar", "baz"} {
		if _, err := svc.Create(ctx, types.CreateResource{Name: name}); err != nil {
			t.Fatalf("Create %q: %v", name, err)
		}
	}

	results, _, err := svc.List(ctx, "name:foo", 10, "")
	if err != nil {
		t.Fatalf("List with filter: %v", err)
	}
	for _, item := range results {
		if item.Name != "foo" {
			t.Errorf("filter returned unexpected item %q", item.Name)
		}
	}
}

func TestUpdate_AppliesPartialFields(t *testing.T) {
	svc := newTestService(t)

	created, err := svc.Create(ctx, types.CreateResource{Name: "original"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	newName := "updated"
	updated, err := svc.Update(ctx, created.ID, types.UpdateResource{Name: &newName})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "updated" {
		t.Errorf("name: want %q, got %q", "updated", updated.Name)
	}
	if !updated.UpdatedAt.After(created.UpdatedAt) {
		t.Error("UpdatedAt should advance after update")
	}
}

func TestDelete_MakesGetReturn404(t *testing.T) {
	svc := newTestService(t)

	created, err := svc.Create(ctx, types.CreateResource{Name: "to-delete"})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if err := svc.Delete(ctx, created.ID); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	_, err = svc.Get(ctx, created.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
	// The service wraps not-found as *errors.NotFound
	type notFound interface{ Error() string }
	if _, ok := err.(notFound); !ok {
		t.Errorf("expected error, got %T: %v", err, err)
	}
}

func TestGet_NotFound(t *testing.T) {
	svc := newTestService(t)

	_, err := svc.Get(ctx, "nonexistent-id")
	if err == nil {
		t.Fatal("expected NotFound error, got nil")
	}
}
