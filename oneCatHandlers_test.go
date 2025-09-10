package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// --- Helpers ---

func setDB(t *testing.T, db map[string]Cat) {
	t.Helper()
	orig := catsDatabase
	catsDatabase = db
	t.Cleanup(func() { catsDatabase = orig })
}

func newReq(method, path, body string) *http.Request {
	r := httptest.NewRequest(method, path, strings.NewReader(body))
	return r
}

func stringSliceAsSet(xs []string) map[string]struct{} {
	s := make(map[string]struct{}, len(xs))
	for _, x := range xs {
		s[x] = struct{}{}
	}
	return s
}

// --- listMapKeys ---

func TestListMapKeys_Empty(t *testing.T) {
	got := listMapKeys(map[string]Cat{})
	if len(got) != 0 {
		t.Fatalf("want empty slice, got %v", got)
	}
}

func TestListMapKeys_AllKeys(t *testing.T) {
	m := map[string]Cat{
		"a": {Name: "A"},
		"b": {Name: "B"},
		"c": {Name: "C"},
	}
	got := listMapKeys(m)
	setGot := stringSliceAsSet(got)

	// On vérifie que TOUTES les clés sont là (l'ordre n'est pas garanti).
	for k := range m {
		if _, ok := setGot[k]; !ok {
			t.Fatalf("missing key %q in %v", k, got)
		}
	}
	if len(got) != len(m) {
		t.Fatalf("got %d keys, want %d", len(got), len(m))
	}
}

// --- listCats ---

func TestListCats_ReturnsDBKeys(t *testing.T) {
	setDB(t, map[string]Cat{
		"id1": {Name: "Toto"},
		"id2": {Name: "Mimi"},
	})

	status, resp := listCats(newReq("GET", "/cats", ""))

	if status != http.StatusOK {
		t.Fatalf("status=%d, want %d", status, http.StatusOK)
	}
	keys, ok := resp.([]string)
	if !ok {
		t.Fatalf("resp type = %T, want []string", resp)
	}

	gotSet := stringSliceAsSet(keys)
	for _, k := range []string{"id1", "id2"} {
		if _, ok := gotSet[k]; !ok {
			t.Fatalf("missing key %q in %v", k, keys)
		}
	}
	if len(keys) != 2 {
		t.Fatalf("got %d keys, want 2", len(keys))
	}
}

// --- createCat ---

func TestCreateCat_Success(t *testing.T) {
	setDB(t, map[string]Cat{}) // base propre

	body := `{"name":"Mimi","color":"Black","birthDate":"2024-01-01"}`
	status, resp := createCat(newReq("POST", "/cats", body))

	if status != http.StatusCreated {
		t.Fatalf("status=%d, want %d", status, http.StatusCreated)
	}

	// L'API renvoie l'ID (string). On vérifie qu'il s'agit d'un UUID valide.
	id, ok := resp.(string)
	if !ok {
		t.Fatalf("resp type = %T, want string (cat id)", resp)
	}
	if _, err := uuid.Parse(id); err != nil {
		t.Fatalf("returned id is not a valid UUID: %q (%v)", id, err)
	}

	// Vérifier que le chat est bien enregistré avec cet ID et les bons champs.
	cat, ok := catsDatabase[id]
	if !ok {
		t.Fatalf("cat with id %q not found in DB", id)
	}
	if cat.ID != id || cat.Name != "Mimi" || cat.Color != "Black" || cat.BirthDate != "2024-01-01" {
		t.Fatalf("saved cat = %#v, fields mismatch", cat)
	}
}

func TestCreateCat_BadJSON(t *testing.T) {
	setDB(t, map[string]Cat{"keep": {Name: "Keep"}})
	before := len(catsDatabase)

	status, resp := createCat(newReq("POST", "/cats", "{not-json"))
	if status != http.StatusBadRequest {
		t.Fatalf("status=%d, want %d", status, http.StatusBadRequest)
	}
	if resp != "Invalid JSON input" {
		t.Fatalf("resp=%v, want %q", resp, "Invalid JSON input")
	}
	// DB ne doit pas être modifiée.
	if len(catsDatabase) != before {
		t.Fatalf("DB changed on bad JSON: before=%d after=%d", before, len(catsDatabase))
	}
}

func TestCreateCat_IgnoresProvidedID(t *testing.T) {
	setDB(t, map[string]Cat{})

	// Même si le client fournit un id, il doit être remplacé par un UUID généré.
	body := `{"id":"hacker","name":"Neo"}`
	status, resp := createCat(newReq("POST", "/cats", body))
	if status != http.StatusCreated {
		t.Fatalf("status=%d, want %d", status, http.StatusCreated)
	}
	id := resp.(string)
	if id == "hacker" {
		t.Fatalf("server should not accept client-provided id")
	}
	cat := catsDatabase[id]
	if cat.ID != id || cat.Name != "Neo" {
		t.Fatalf("saved cat mismatch: %#v", cat)
	}
}
