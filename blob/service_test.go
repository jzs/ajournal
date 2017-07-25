package blob_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"strings"
	"testing"
	"time"

	"github.com/sketchground/ajournal/blob"
)

func TestCreate(t *testing.T) {
	mr := &metaRepo{map[string]*blob.File{}}
	s := blob.NewService(mr)

	key := "journal/1/images"
	mimetype := "text/plain"
	data := bytes.NewReader([]byte("Hello world!"))
	details, err := s.Create(key, mimetype, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if details.Key != key {
		t.Errorf("Expected key %v, got %v", key, details.Key)
	}
}

func TestDetails(t *testing.T) {
	mr := &metaRepo{map[string]*blob.File{}}
	s := blob.NewService(mr)

	key := "journal/1/images"
	mimetype := "text/plain"
	data := bytes.NewReader([]byte("Hello world!"))
	_, err := s.Create(key, mimetype, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	details, err := s.Details(key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if details.Key != key {
		t.Errorf("Expected key %v, got %v", key, details.Key)
	}
}

func TestValue(t *testing.T) {
	mr := &metaRepo{map[string]*blob.File{}}
	s := blob.NewService(mr)

	key := "journal/1/images"
	mimetype := "text/plain"
	data := bytes.NewReader([]byte("Hello world!"))
	_, err := s.Create(key, mimetype, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	val, err := s.Value(key)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	ndata, err := ioutil.ReadAll(val)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if string(ndata) != "Hello world!" {
		t.Fatalf("Expected 'Hello world!', got %v", string(ndata))
	}

}

func TestList(t *testing.T) {
	mr := &metaRepo{map[string]*blob.File{}}
	s := blob.NewService(mr)

	key := "journal/1/images"
	mimetype := "text/plain"
	data := bytes.NewReader([]byte("Hello world!"))
	_, err := s.Create(key, mimetype, data)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	files, err := s.List("journal/1")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected one file, got %v", len(files))
	}
	if files[0].Key != key {
		t.Errorf("Expected key %v, got %v", key, files[0].Key)
	}
}

type metaRepo struct {
	data map[string]*blob.File
}

func (m *metaRepo) Get(key string) (*blob.File, error) {
	if details, ok := m.data[key]; ok {
		return details, nil
	}
	return nil, errors.New("Does not exist")
}

func (m *metaRepo) Head(key string) (*blob.File, error) {
	if details, ok := m.data[key]; ok {
		return details, nil
	}
	return nil, errors.New("Does not exist")
}

func (m *metaRepo) Create(key, mimetype string, r io.Reader) (*blob.File, error) {
	f := &blob.File{
		MIMEType: mimetype,
		Key:      key,
		Created:  time.Now(),
		Reader:   r,
	}
	m.data[key] = f
	return f, nil
}

func (m *metaRepo) List(keypath string) ([]*blob.File, error) {
	files := []*blob.File{}
	for k, v := range m.data {
		if strings.Contains(k, keypath) {
			files = append(files, v)
		}
	}
	return files, nil
}
