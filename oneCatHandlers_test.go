package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// petite aide pour préparer la "DB" mémoire
func setDB(t *testing.T, db map[string]Cat) {
	t.Helper()
	orig := catsDatabase
	catsDatabase = db
	t.Cleanup(func() { catsDatabase = orig })
}

// Test très simple : getCat retourne 404 quand l'ID n'existe pas
func TestGetCat_NotFound(t *testing.T) {
	setDB(t, map[string]Cat{}) // DB vide

	r := httptest.NewRequest("GET", "/cats/nope", nil)
	r.SetPathValue("catId", "nope") // Go 1.25 : injecte le paramètre de path

	status, resp := getCat(r)
	if status != http.StatusNotFound {
		t.Fatalf("status=%d, want %d", status, http.StatusNotFound)
	}
	if resp != "Cat not found" {
		t.Fatalf("resp=%v, want %q", resp, "Cat not found")
	}
}
