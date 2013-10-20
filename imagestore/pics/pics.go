// Package pics contains the internal API for accessing images.
package pics

import (
	"fmt"
	"regexp"
	"strings"

	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"appengine/image"
)

const (
	Kind = "Pic"
)

var rxSHA1 = regexp.MustCompile("^[0-9a-f]{40}$")

type Pic struct {
	ID string `datastore:"-"`

	Key appengine.BlobKey `datastore:"key,noindex"`
	URL string            `datastore:"url,noindex"`
}

// Get gets a pic from the datastore.
func Get(c appengine.Context, id string) (*Pic, error) {
	var p Pic
	k := NewKey(c, strings.ToLower(id))
	err := datastore.Get(c, k, &p)
	switch err {
	default:
		return nil, err
	case datastore.ErrNoSuchEntity:
		return nil, nil
	case nil:
		p.ID = k.StringID()
		return &p, nil
	}
}

// Create creates a new pic and stores it in the datastore.
func Create(c appengine.Context, b *blobstore.BlobInfo) (*Pic, error) {
	id := strings.ToLower(strings.SplitN(b.Filename, ".", 2)[0])
	if rxSHA1.FindString(id) == "" {
		return nil, fmt.Errorf("%s is not a sha1 hash", id)
	}
	p := &Pic{
		ID:  id,
		Key: b.BlobKey,
	}
	url, err := image.ServingURL(c, p.Key, nil)
	if err != nil {
		c.Errorf("Failed to create serving URL (%v).", err)
		return nil, err
	}
	p.URL = url.String()
	_, err = datastore.Put(c, NewKey(c, id), p)
	return p, err
}

// NewKey creates a new pic key.
func NewKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, Kind, id, 0, nil)
}
