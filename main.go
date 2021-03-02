package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// template for html
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// handle http requests
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// compile template once ( lazy initialisation )
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, nil)
}

// start of the program
func main() {

	// serve "/" root

	// pass address of newly created object of type templateHandler
	http.Handle("/", &templateHandler{filename: "skychat.html"})

	// start serving ...
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}

}