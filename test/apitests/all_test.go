package apitests

import (
	"fmt"
	"net/http"
	"testing"
)

var initCatId string

func init() {
	ids := []string{}
	_ = call("GET", "/cats", nil, nil, &ids)

	// 2) Supprimer tout
	for _, id := range ids {
		code := 0
		_ = call("DELETE", "/cats/"+id, nil, &code, nil)
		fmt.Println("DELETE /cats/"+id, "->", code)
	}

	_ = call("POST", "/cats", &CatModel{Name: "Toto"}, nil, &initCatId)
	fmt.Println("Created cat:", initCatId)
}

func TestGetCats(t *testing.T) {
	code := 0
	result := []string{}

	if err := call("GET", "/cats", nil, &code, &result); err != nil {
		t.Fatalf("Request error: %v", err)
	}

	fmt.Println("GET /cats ->", code, result)

	if code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}

	if len(result) != 1 {
		t.Fatalf("expected 1 id, got %d (%v)", len(result), result)
	}

	if result[0] != initCatId {
		t.Fatalf("expected id %s, got %v", initCatId, result)
	}
}
