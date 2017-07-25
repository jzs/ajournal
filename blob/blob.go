package blob

import (
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
	Filename string
	MIMEType string
	Created  time.Time
	Reader   io.Reader
	Links    Links
}

// Repository describes a repository for fetching the blob data
type Repository interface {
	Get(key string) (*File, error)
	Create(key, mimetype string, r io.Reader) (*File, error)
	Head(key string) (*File, error)
	List(keypath string) ([]*File, error)
}
