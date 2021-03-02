package main

import (
	"flag"
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

	// pass request data ( which includes the host address )
	t.templ.Execute(w, r)
}

// start of the program
func main() {

	// parse command-line arguments ( default : 8080 )
	var addr = flag.String("addr", ":8080", "Address for the app")
	flag.Parse()

	// create a room
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "skychat.html"})
	http.Handle("/room", r)

	// start the room ( in separate go routine )
	go r.run()

	// start the webserver ( main routine )
	log.Println("Starting web server at : ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalln("ListenAndServe:", err)
	}

}
