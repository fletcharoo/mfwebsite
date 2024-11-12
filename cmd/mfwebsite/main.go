package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/fletcharoo/snest"
)

type Config struct {
	Port string `snest:"PORT"`
}

//go:embed style.css
var styleCSS string
var workingDir string

func main() {
	var conf Config
	if err := snest.Load(&conf); err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	http.HandleFunc("/style.css", handlerFactory(styleCSS))
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	addMarkdownRoutes(workingDir)

	addr := ":" + conf.Port
	log.Println("Serving on", addr)
	http.ListenAndServe(addr, nil)
}

func handlerFactory(data string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, data)
	}
}

func addMarkdownRoutes(dir string) (err error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		err = fmt.Errorf("failed to read dir %q: %w", dir, err)
		return
	}

	for _, entry := range entries {
		entryName := entry.Name()
		pathFull := filepath.Join(dir, entryName)
		if entry.IsDir() {
			if err = addMarkdownRoutes(pathFull); err != nil {
				return
			}
		}

		if strings.ToLower(filepath.Ext(pathFull)) != ".md" {
			continue
		}

		var contents []byte
		contents, err = os.ReadFile(pathFull)
		if err != nil {
			err = fmt.Errorf("failed to read file %q: %w", pathFull, err)
			return
		}

		entryNameNoExt := strings.Split(entryName, ".")[0]
		pathSplit := strings.Split(pathFull, string(os.PathSeparator))
		if len(pathSplit) != 0 {
			pathSplit = pathSplit[2 : len(pathSplit)-1]
		}

		pathSplit = append(pathSplit, entryNameNoExt)
		addr := "/" + strings.Join(pathSplit, "/")
		if addr == "/index" {
			addr = "/"
		}

		http.HandleFunc(addr, handlerFactory(string(contents)))
		log.Println("Registered", addr)
	}

	return nil
}
