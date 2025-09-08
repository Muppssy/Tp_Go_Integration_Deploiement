package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type Cat struct {
	Name      string `json:"name"`
	ID        string `json:"id,omitempty"`
	BirthDate string `json:"birthDate,omitempty"`
	Color     string `json:"color,omitempty"`
}

// Simple in-memory database, for demo purpose
var catsDatabase = map[string]Cat{
	"id1": {Name: "Toto", Color: "Grey", BirthDate: "2023-04-16"},
}

func listMapKeys(aMap map[string]Cat) []string {
	results := []string{}

	for catID := range aMap {
		results = append(results, catID)
	}

	return results
}

func listCats(req *http.Request) (int, any) {
	Logger.Info("Listing the cats")
	return http.StatusOK, listMapKeys(catsDatabase)
}

func createCat(req *http.Request) (int, any) {

	// Decode the request body into a Cat structure
	decoder := json.NewDecoder(req.Body)
	var catCreationData Cat
	err := decoder.Decode(&catCreationData)
	if err != nil {
		Logger.Info("Unable to parse the JSON input for cat creation")
		return http.StatusBadRequest, "Invalid JSON input"
	}

	Logger.Info("Creating the cat: ", catCreationData)

	// Creating the new cat's ID and storing the Cat
	newCatID := uuid.New().String()
	catCreationData.ID = newCatID

	catsDatabase[newCatID] = catCreationData

	Logger.Infof("Cat '%s' saved into the DB", newCatID)
	return http.StatusCreated, newCatID
}
