package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"

	"github.com/fletcharoo/snest"
)

type Config struct {
	Port string `snest:"PORT"`
}

//go:embed style.css
var styleCSS string

func main() {
	var conf Config
	if err := snest.Load(&conf); err != nil {
		log.Fatalf("failed to load config: %s", err)
	}

	http.HandleFunc("/style.css", styleCSSHandler)
	addr := ":" + conf.Port
	log.Printf("Serving on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

func styleCSSHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Fprintf(w, styleCSS)
}
