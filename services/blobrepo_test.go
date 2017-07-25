package services_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/sketchground/ajournal/services"
)

func TestS3Create(t *testing.T) {
	endpoint := os.Getenv("AJ_S3_ENDPOINT")
	accessKey := os.Getenv("AJ_S3_ACCESSKEY")
	secretKey := os.Getenv("AJ_S3_SECRETKEY")
	if endpoint == "" || accessKey == "" || secretKey == "" {
		t.Fatalf("S3 Credentials not specified, cannot test s3 integration")
	}

	repo := services.NewS3Repo(endpoint, accessKey, secretKey, "ajournal-test")
	_, err := repo.Create("test/test.md", "text/plain", bytes.NewReader([]byte("Hello world!")))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Get file see if we can fetch it, and head file...
	file, err := repo.Get("test/test.md")
	if err != nil {
		t.Fatalf("Expected go fetch file, got %v", err)
	}
	if file.Key != "test/test.md" {
		t.Errorf("Expected key test/test.md, got: %v", file.Key)
	}
	if file.MIMEType != "text/plain" {
		t.Errorf("Expected mimetype text/plain , got: %v", file.MIMEType)
	}

	files, err := repo.List("test/")
	if err != nil {
		t.Fatalf("Expected go fetch file list, got %v", err)
	}
	if len(files) != 1 {
		t.Fatalf("Expected 1 file, got %v", len(files))
	}
	if files[0].Key != "test/test.md" {
		t.Errorf("Expected key test/test.md, got: %v", file.Key)
	}
}
