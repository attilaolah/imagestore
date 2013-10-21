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

var rxSHA1 = regexp.MustCompile("(?i)^[0-9a-f]{40}\\.jpe?g$")

type Pic struct {
	ID string `datastore:"-"`

	Key appengine.BlobKey `datastore:"key,noindex"`
	URL string            `datastore:"url,noindex"`
}

// Get gets an image from the datastore.
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

// Create creates a new image and stores it in the datastore.
func Create(c appengine.Context, b *blobstore.BlobInfo) (*Pic, error) {
	if rxSHA1.FindString(b.Filename) == "" {
		if err := blobstore.Delete(c, b.BlobKey); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s does not match (?i)^[0-9a-f]{40}\\.jpe?g$", b.Filename)
	}
	id := strings.ToLower(strings.SplitN(b.Filename, ".", 2)[0])
	// Check for an existing image first
	p, err := Get(c, id)
	if err != nil {
		return nil, err
	}
	if p != nil {
		if err = p.Delete(c); err != nil {
			return nil, err
		}
	}
	// Create the new image
	p = &Pic{
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

// Delete removes an image from both the blobstore and the datastore.
func (p *Pic) Delete(c appengine.Context) (err error) {
	if err = blobstore.Delete(c, p.Key); err == nil {
		err = datastore.Delete(c, NewKey(c, p.ID))
	}
	return
}

// NewKey creates a new image key.
func NewKey(c appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(c, Kind, id, 0, nil)
}
