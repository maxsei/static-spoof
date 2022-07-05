package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Get the filepath to an image to serve in the index page.
	if len(os.Args) != 2 {
		fmt.Println("must supply a filepath to an image")
		os.Exit(1)
	}

	// Filepath to resource.
	fp := "./static/matrix-reloaded.jpg"
	// Check if path is valid and isn't a directory.
	fi, err := os.Stat(fp)
	if os.IsNotExist(err) {
		log.Fatal(err)
	}
	if fi.IsDir() {
		log.Fatal(fmt.Errorf("%s must be a file", fp))
	}

	// Load index html template.
	tmplMap := map[string]string{
		"Filename": fi.Name(),
		"Filepath": fp,
	}
	tmplBuf, err := ExecuteAndParseTemplateUnbuffered("./static/index.tmpl", tmplMap)
	if err != nil {
		log.Fatal(err)
	}
	// Index html handler.
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		n, err := w.Write(tmplBuf)
		if err != nil {
			log.Print(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", fmt.Sprintf("%v", n))
	})

	// Path mappings are the name of the file and an arbitrary exenstion
	// (different-path) of the path path.
	var pathMappings = []*PathMapping{
		// Exposed                                     // Actual
		{fmt.Sprintf("/%s", fi.Name()), fp},
		{fmt.Sprintf("/different-path/%s", fi.Name()), fp},
	}
	for _, p := range pathMappings {
		http.Handle(p.Exposed, p)
	}

	// Listen and serve at port.
	const port int = 8080
	fmt.Printf("listen and serve at: localhost:%d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

type PathMapping struct {
	// Exposed is the path that is exposed to the user of this API/Service
	Exposed string
	// Actual is the real filepath of the resource.
	Actual string
}

func (p *PathMapping) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	// Open file and get file info.
	f, err := os.Open(p.Actual)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	stat, err := f.Stat()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Write response and headers.
	n, err := io.Copy(w, f)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%v", n))
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Date", fmt.Sprintf("%v", time.Now()))
	w.Header().Set("Last-Modified", fmt.Sprintf("%v", stat.ModTime()))
}


func ExecuteAndParseTemplateUnbuffered(path string, data interface{}) ([]byte, error) {
	tmpl, err := template.ParseFiles("./static/index.tmpl")
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer([]byte{})
	if err := tmpl.Execute(buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
