package blob

import (
	"fmt"
	"io"
	"time"
)

// Links presents links for the files
type Links struct {
	Orig  string
	Thumb string
}

// File represents details about an uploaded blob
type File struct {
	Key      string
	Filename string    `json:"Filename,omitempty"`
	MIMEType string    `json:"MIMEType,omitempty"`
	Created  time.Time `json:"Created,omitempty"`
	Reader   io.Reader `json:"Reader,omitempty"`
	Links    Links     `json:"Links,omitempty"`
}

// Repository describes a repository for fetching the blob data
type Repository interface {
	Get(key string) (*File, error)
	Create(key, mimetype string, r io.Reader) (*File, error)
	Head(key string) (*File, error)
	List(keypath string) ([]*File, error)
}

// FileFromKey generates a file struct based on a given key.
func FileFromKey(key string) File {
	// Generate urls etc...
	return File{
		Key: key,
		Links: Links{
			Orig: fmt.Sprintf("/api/blobs/%v", key),
		},
	}
}
