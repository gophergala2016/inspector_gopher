package main

import (
	"log"
	"net/http"
	"fmt"
	"github.com/gophergala2016/inspector_gopher"
	"flag"
)

var webRoot = flag.String("webroot", "public", "Relative or absolute path to the directory where the static servable files are stored.")

func main() {
	flag.Parse()

	fs := http.FileServer(http.Dir(*webRoot))

	http.Handle("/", fs)

	fmt.Println("Starting to serve!")
		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		repoUrl := r.URL.Query().Get("repo")
		coordinator := inspector.NewCoordinator(repoUrl)

		fmt.Fprintf(w, "%q", coordinator.Heatmap())
	})

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
