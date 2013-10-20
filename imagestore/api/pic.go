package api

import (
	"net/http"
	"strings"

	"appengine"

	"github.com/gorilla/mux"

	"imagestore/pics"
)

func Pic(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	for _, a := range strings.Split(r.Header.Get("Accept"), ",") {
		a = strings.SplitN(a, ";", 2)[0]
		switch a {
		case "image/jpeg", "*/*", "":
			redirect(w, r)
			return
		}
	}
	http.Error(w, "", http.StatusNotAcceptable)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	id := mux.Vars(r)["id"]
	p, err := pics.Get(c, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if p == nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, p.URL, http.StatusMovedPermanently)
}
