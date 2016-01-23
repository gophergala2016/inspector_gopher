package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gophergala2016/inspector_gopher"
)

func main() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		repoUrl := r.URL.Query().Get("repo")
		coordinator := inspector.NewCoordinator(repoUrl)

		fmt.Fprintf(w, "%q", coordinator.Heatmap())
	})

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
