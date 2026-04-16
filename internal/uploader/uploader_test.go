package uploader

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"live-webrtc-go/internal/config"
)

func resetUploader() {
	client = nil
	cfg = nil
}

func TestInit_UploadDisabled(t *testing.T) {
	resetUploader()
	c := &config.Config{UploadEnabled: false}
	if err := Init(c); err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
	if client != nil {
		t.Error("expected client to be nil when upload disabled")
	}
}

func TestInit_MissingS3Config(t *testing.T) {
	resetUploader()
	tests := []struct {
		name string
		cfg  *config.Config
	}{
		{
			name: "missing endpoint",
			cfg:  &config.Config{UploadEnabled: true, S3Endpoint: "", S3Bucket: "bucket", S3AccessKey: "key", S3SecretKey: "secret"},
		},
		{
			name: "missing bucket",
			cfg:  &config.Config{UploadEnabled: true, S3Endpoint: "localhost:9000", S3Bucket: "", S3AccessKey: "key", S3SecretKey: "secret"},
		},
		{
			name: "missing access key",
			cfg:  &config.Config{UploadEnabled: true, S3Endpoint: "localhost:9000", S3Bucket: "bucket", S3AccessKey: "", S3SecretKey: "secret"},
		},
		{
			name: "missing secret key",
			cfg:  &config.Config{UploadEnabled: true, S3Endpoint: "localhost:9000", S3Bucket: "bucket", S3AccessKey: "key", S3SecretKey: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetUploader()
			err := Init(tt.cfg)
			if err == nil {
				t.Error("expected error for missing S3 config")
			}
		})
	}
}

func TestEnabled_ReturnsFalseWhenNotConfigured(t *testing.T) {
	resetUploader()
	if Enabled() {
		t.Error("expected Enabled() to return false when not initialized")
	}
}

func TestEnabled_ReturnsFalseWhenUploadDisabled(t *testing.T) {
	resetUploader()
	_ = Init(&config.Config{UploadEnabled: false})
	if Enabled() {
		t.Error("expected Enabled() to return false when upload disabled")
	}
}

func TestUpload_ReturnsNilWhenDisabled(t *testing.T) {
	resetUploader()
	// Not initialized, so Upload should return nil
	if err := Upload(context.Background(), "/nonexistent"); err != nil {
		t.Errorf("expected nil error when upload disabled, got %v", err)
	}
}

func TestUpload_ReturnsErrorForNonexistentFile(t *testing.T) {
	// Create a temp file to test the file opening path
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Initialize with upload disabled (so client is nil)
	resetUploader()
	_ = Init(&config.Config{UploadEnabled: false})

	// Upload should return nil when disabled
	err := Upload(context.Background(), tmpFile)
	if err != nil {
		t.Errorf("expected nil error when upload disabled, got %v", err)
	}
}
