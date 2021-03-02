package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// implements ServeHTTP method
type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// look for cookies
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie {
		// if not found i.e not authorised
		// redirect to login page
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// else if auth succeeds ~ call next handler
	h.next.ServeHTTP(w, r)
}

// MustAuth creates authHandler that wraps any other handler
// ( to easily add auth to our code )
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}

}

// loginHandler handles the third party login process
// format: /auth/{action}/{provider

func loginHandler(w http.ResponseWriter, r *http.Request) {

	// break path to pull action & provider
	segs := strings.Split(r.URL.Path, "/")

	// [ To Do ] : catch garbage cases

	action := segs[2]
	provider := segs[3]

	switch action {
	// if action value is known , run specific code
	case "login":
		log.Println("TODO handle login for", provider)
	// else return error
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action is %s not supported", action)
	}

}
