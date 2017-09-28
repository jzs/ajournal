package services

import (
	"io"
	"net/http"
	"strings"

	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/sketchground/ajournal/blob"
	"github.com/sketchground/ajournal/utils"
)

// NewS3Repo returns a new blob repository using an underlying s3 service
func NewS3Repo(endpoint, accessKey, secretKey, bucket string) (*S3Repo, error) {
	ssl := true
	if strings.Contains(endpoint, "127.0.0.1") {
		ssl = false
	}
	client, err := minio.New(endpoint, accessKey, secretKey, ssl)
	if err != nil {
		return nil, errors.Wrap(err, "Failed configuring minio s3 service")
	}
	return &S3Repo{
		client: client,
		bucket: bucket,
	}, nil
}

// S3Repo struct
type S3Repo struct {
	client *minio.Client
	bucket string
}

// Get gets a file from s3
func (m *S3Repo) Get(key string) (*blob.File, error) {
	object, err := m.client.GetObject(m.bucket, key)
	if err != nil {
		return nil, err
	}

	stat, err := object.Stat()
	if err != nil {
		erresp := minio.ToErrorResponse(err)
		if erresp.Code == "NoSuchKey" {
			return nil, utils.NewAPIError(err, http.StatusNotFound, "Blob not found")
		}
		return nil, err
	}

	return &blob.File{
		Key:      stat.Key,
		MIMEType: stat.ContentType,
		//Created:  stat.Created,
		Reader: object,
	}, nil
}

// Head returns stats about a file
func (m *S3Repo) Head(key string) (*blob.File, error) {
	object, err := m.client.GetObject(m.bucket, key)
	if err != nil {
		return nil, err
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, err
	}

	return &blob.File{
		Key:      stat.Key,
		MIMEType: stat.ContentType,
		//Created:  stat.Created,
	}, nil
}

// Create creates a new blob in the storage
func (m *S3Repo) Create(key, mimetype string, r io.Reader) (*blob.File, error) {
	_, err := m.client.PutObject(m.bucket, key, r, mimetype)
	if err != nil {
		return nil, err
	}
	return &blob.File{}, nil
}

// List lists files in a bucket with the given prefix
func (m *S3Repo) List(keypath string) ([]*blob.File, error) {
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectCh := m.client.ListObjectsV2(m.bucket, keypath, false, doneCh)
	files := []*blob.File{}
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		files = append(files, &blob.File{Key: object.Key})
	}

	return files, nil
}

// CreateBucket creates a new bucket in s3
func (m *S3Repo) CreateBucket(name, location string) error {
	err := m.client.MakeBucket(name, location)
	if err != nil {
		return err
	}
	return nil
}
