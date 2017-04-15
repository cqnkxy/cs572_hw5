package main

import (
	"log"
	"net/http"

	"views"
)

const (
	PORT = ":8080"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("./templates")))
	http.Handle("/static/", http.StripPrefix(
		"/static", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/search/", views.Search)
	http.HandleFunc("/correct/", views.Correct)
	http.HandleFunc("/suggest/", views.Suggest)
	log.Printf("Server running on localhost%s\n", PORT)
	http.ListenAndServe(PORT, nil)
}
