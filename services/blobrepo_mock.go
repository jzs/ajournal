package services

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sketchground/ajournal/blob"
	"github.com/sketchground/ajournal/utils"
)

// NewS3MockRepo returns a new mock blob repository
func NewS3MockRepo() *S3MockRepo {
	return &S3MockRepo{
		data: map[string]*blob.File{},
	}
}

// S3MockRepo struct
type S3MockRepo struct {
	data map[string]*blob.File
}

// Get gets a file from s3
func (m *S3MockRepo) Get(key string) (*blob.File, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, utils.NewAPIError(nil, http.StatusNotFound, "Blob not found")
	}

	return val, nil
}

// Head returns stats about a file
func (m *S3MockRepo) Head(key string) (*blob.File, error) {
	val, ok := m.data[key]
	if !ok {
		return nil, utils.NewAPIError(nil, http.StatusNotFound, "Blob not found")
	}

	return val, nil
}

// Create creates a new blob in the storage
func (m *S3MockRepo) Create(key, mimetype string, r io.Reader) (*blob.File, error) {
	m.data[key] = &blob.File{
		Key:      key,
		MIMEType: mimetype,
		Reader:   r,
	}
	return m.data[key], nil
}

// List lists files in a bucket with the given prefix
func (m *S3MockRepo) List(keypath string) ([]*blob.File, error) {
	return nil, errors.New("TODO: Implement me")
}

// CreateBucket creates a new bucket in s3
func (m *S3MockRepo) CreateBucket(name, location string) error {
	return errors.New("TODO: Implement me")
}
