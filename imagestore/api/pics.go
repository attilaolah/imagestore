package api

import (
	"encoding/json"
	"net/http"

	"appengine"
	"appengine/blobstore"

	"imagestore/pics"
)

func Pics(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	blobs, _, err := blobstore.ParseUpload(r)
	if err != nil {
		// No submitted blobs, create a new upload session
		url, err := blobstore.UploadURL(c, r.URL.Path, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url.String(), http.StatusFound)
		return
	}
	files := blobs["file"]
	if len(files) == 0 {
		http.Error(w, "no files submitted", http.StatusBadRequest)
		return
	}
	urls := make([]string, len(files))
	for i, f := range files {
		p, err := pics.Create(c, f)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		urls[i] = p.URL
	}
	w.WriteHeader(http.StatusAccepted)
	if err := json.NewEncoder(w).Encode(urls); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
