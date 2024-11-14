package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"

	"github.com/fletcharoo/snest"
)

type Config struct {
	Port string `snest:"PORT"`
}

//go:embed style.css
var styleCSS string
var workingDir string

func main() {
	// Load service configs.
	var conf Config
	if err := snest.Load(&conf); err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %s", err)
	}

	// Add API routes.
	http.HandleFunc("/style.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		fmt.Fprintf(w, styleCSS)
	})

	addMarkdownRoutes(workingDir)

	// Start service.
	addr := ":" + conf.Port
	log.Println("Serving on", addr)
	http.ListenAndServe(addr, nil)
}

// handlerFactory returns a function that implements the input for
// http.HandleFunc which writes the provided string to the response writer.
func handlerFactory(data string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, data)
	}
}

// addMarkdownRoutes finds all markdown files within the current working
// directory and all sub directories, renders the markdown as HTML, and adds
// routes to return the generated HTML.
// The ".md" extension is stripped when creating the route.
// Any file in the working directory called "index.md" is inferred to be the
// root route.
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

		http.HandleFunc(addr, handlerFactory(mdToHTML(contents)))
		log.Println("Registered", addr)
	}

	return nil
}

// mdToHTML accepts markdown and renders it as a HTML page.
func mdToHTML(md []byte) (renderedHTML string) {
	// Create markdown parser with extensions.
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// Create HTML renderer with extensions.
	htmlFlags := html.CommonFlags | html.HrefTargetBlank | html.CompletePage
	opts := html.RendererOptions{
		Flags: htmlFlags,
		CSS:   "style.css",
		Head:  []byte(`<meta name="viewport" content="width=device-width, initial-scale=1">`),
	}
	renderer := html.NewRenderer(opts)
	renderedHTML = string(markdown.Render(doc, renderer))

	return renderedHTML
}
