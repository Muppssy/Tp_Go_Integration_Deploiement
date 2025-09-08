package main

import "net/http"

func getCat(req *http.Request) (int, any) {
	catID := req.PathValue("catId")
	Logger.Info("Getting the cat: ", catID)

	if cat, found := catsDatabase[catID]; found {
		Logger.Info("Cat found")
		return http.StatusOK, cat
	} else {
		Logger.Info("Cat not found")
		return http.StatusNotFound, "Cat not found"
	}
}
