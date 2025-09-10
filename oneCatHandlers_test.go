package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// --- Helpers ---

func setDB(t *testing.T, db map[string]Cat) {
	t.Helper()
	orig := catsDatabase
	catsDatabase = db
	t.Cleanup(func() { catsDatabase = orig })
}

// Injecte catId pour que getCat/deleteCat puissent lire req.PathValue("catId")
func newReqWithCatID(method, id string) *http.Request {
	r := httptest.NewRequest(method, "/cats/"+id, nil)
	r.SetPathValue("catId", id)
	return r
}

// --- Tests getCat ---

func TestGetCat_Found(t *testing.T) {
	setDB(t, map[string]Cat{
		"id1": {ID: "id1", Name: "Mimi", Color: "Black", BirthDate: "2024-01-01"},
	})

	status, resp := getCat(newReqWithCatID("GET", "id1"))

	if status != http.StatusOK {
		t.Fatalf("status=%d, want %d", status, http.StatusOK)
	}
	cat, ok := resp.(Cat)
	if !ok {
		t.Fatalf("resp type = %T, want Cat", resp)
	}
	if cat.ID != "id1" || cat.Name != "Mimi" {
		t.Fatalf("cat mismatch: %#v", cat)
	}
}

func TestGetCat_NotFound(t *testing.T) {
	setDB(t, map[string]Cat{})

	status, resp := getCat(newReqWithCatID("GET", "nope"))

	if status != http.StatusNotFound {
		t.Fatalf("status=%d, want %d", status, http.StatusNotFound)
	}
	if resp != "Cat not found" {
		t.Fatalf("resp=%v, want %q", resp, "Cat not found")
	}
}

// --- Tests deleteCat ---

func TestDeleteCat_RemovesAndReturnsCat(t *testing.T) {
	setDB(t, map[string]Cat{
		"id1": {ID: "id1", Name: "Toto"},
	})

	before := len(catsDatabase)
	status, resp := deleteCat(newReqWithCatID("DELETE", "id1"))

	if status != http.StatusOK {
		t.Fatalf("status=%d, want %d", status, http.StatusOK)
	}
	cat, ok := resp.(Cat)
	if !ok {
		t.Fatalf("resp type = %T, want Cat", resp)
	}
	if cat.ID != "id1" || cat.Name != "Toto" {
		t.Fatalf("returned cat mismatch: %#v", cat)
	}
	if _, exists := catsDatabase["id1"]; exists {
		t.Fatalf("cat id1 should have been deleted from DB")
	}
	if len(catsDatabase) != before-1 {
		t.Fatalf("DB size mismatch: before=%d after=%d", before, len(catsDatabase))
	}
}
