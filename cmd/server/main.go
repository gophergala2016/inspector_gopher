package main

import (
	"net/http"
	"log"
)

func main() {
  fs := http.FileServer(http.Dir("public"))
  http.Handle("/", fs)

  log.Println("Listening...")
  http.ListenAndServe(":80", nil)
}
