// Package imagestore implements image storage on top of App Engine.
// The Images API is used to serve Images that are stored in the Blobstore.
package imagestore

import (
	"net/http"

	"github.com/gorilla/mux"

	"imagestore/api"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/", api.Pics)
	r.HandleFunc("/{id:(?i)[0-9a-z]{40}}.jpg", api.Pic)

	http.Handle("/", r)
}
