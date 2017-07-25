package blob

import (
	"fmt"
	"io"
)

// Service describes the functionality one can do on blobs
type Service interface {
	Create(key, filetype string, r io.Reader) (*File, error)
	Details(key string) (*File, error)
	Value(key string) (io.Reader, error)
	List(keypath string) ([]*File, error)
}

type service struct {
	br Repository
}

// NewService returns a new instance of the blob service
func NewService(repo Repository) Service {
	return &service{br: repo}
}

func (s *service) Create(key, filetype string, r io.Reader) (*File, error) {
	// TODO: Create a thumbnail as well... if blob is an image...
	return s.br.Create(key, filetype, r)
}

func (s *service) Details(key string) (*File, error) {
	f, err := s.br.Get(key)
	if err != nil {
		return nil, err
	}
	f.Links = Links{
		Orig:  fmt.Sprintf("api/%v", key),
		Thumb: fmt.Sprintf("api/%v_thumb", key),
	}
	return f, nil
}

func (s *service) Value(key string) (io.Reader, error) {
	v, err := s.br.Get(key)
	if err != nil {
		return nil, err
	}
	return v.Reader, nil
}

func (s *service) List(keypath string) ([]*File, error) {
	list, err := s.br.List(keypath)
	if err != nil {
		return nil, err
	}
	for _, f := range list {
		f.Links = Links{
			Orig:  fmt.Sprintf("api/%v", f.Key),
			Thumb: fmt.Sprintf("api/%v_thumb", f.Key),
		}
	}
	return list, nil
}
