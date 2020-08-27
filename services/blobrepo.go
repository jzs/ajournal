package services

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/jzs/ajournal/blob"
	"github.com/jzs/ajournal/utils"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pkg/errors"
)

// NewS3Repo returns a new blob repository using an underlying s3 service
func NewS3Repo(endpoint, accessKey, secretKey, bucket string) (*S3Repo, error) {
	ssl := true
	if strings.Contains(endpoint, "127.0.0.1") {
		ssl = false
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: ssl,
	})
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
	object, err := m.client.GetObject(context.TODO(), m.bucket, key, minio.GetObjectOptions{})
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
	object, err := m.client.GetObject(context.TODO(), m.bucket, key, minio.StatObjectOptions{})
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
	_, err := m.client.PutObject(context.TODO(), m.bucket, key, r, -1, minio.PutObjectOptions{ContentType: mimetype})
	if err != nil {
		return nil, err
	}
	return &blob.File{
		Key:      key,
		MIMEType: mimetype,
	}, nil
}

// List lists files in a bucket with the given prefix
func (m *S3Repo) List(keypath string) ([]*blob.File, error) {
	objectCh := m.client.ListObjects(context.TODO(), m.bucket, minio.ListObjectsOptions{Prefix: keypath})
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
	err := m.client.MakeBucket(context.TODO(), name, minio.MakeBucketOptions{Region: location})
	if err != nil {
		return err
	}
	return nil
}
