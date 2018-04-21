package main

import (
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/search", handleSearch)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSearch(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")
	if query == ""{
		http.Error(w, "no query", http.StatusBadRequest)
		return 
	}
}