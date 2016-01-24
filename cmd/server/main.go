package main

import (
	"flag"
	"fmt"
	"github.com/gophergala2016/inspector_gopher"
	"log"
	"net/http"
	"os"
	"strconv"
)

var webRoot = flag.String("webroot", os.Getenv("GOPATH")+string(os.PathSeparator)+"src/github.com/gophergala2016/inspector_gopher/public", "Relative or absolute path to the directory where the static servable files are stored.")
var repoDir = flag.String("repodir", "/tmp", "The directory in which to store cloned repositories.")
func main() {
	flag.Parse()

	inspector.SetTempDir(*repoDir)

	fs := http.FileServer(http.Dir(*webRoot))

	http.Handle("/", fs)

	fmt.Println("Starting to serve!")
	http.HandleFunc("/benchmark", func(w http.ResponseWriter, r *http.Request) {
		repoUrl := r.URL.Query().Get("repo")
		depth := r.URL.Query().Get("depth")

		f, _ := strconv.ParseInt(depth, 0, 0)
		log.Printf("repoName: %s, depth: %d", repoUrl, f)

		time := inspector.HarvestBenched(repoUrl, int(f))

		w.Write([]byte(fmt.Sprintf("%f", time)))
		w.WriteHeader(200)
	})
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		repoUrl := r.URL.Query().Get("repo")
		coordinator := inspector.NewCoordinator(repoUrl)

		fmt.Fprintf(w, "%q", coordinator.Heatmap())
	})

	log.Println("Listening...")
	http.ListenAndServe(":8080", nil)
}
