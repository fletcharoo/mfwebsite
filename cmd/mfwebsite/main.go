package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", response)
	log.Println("Serving on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func response(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, "Hello, world!")
}
