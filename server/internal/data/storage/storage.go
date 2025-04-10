package storage

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/url"
	"server/internal/config"
	"strings"
	"time"
)

type MiniRepo struct {
	client     *minio.Client
	bucketName string
}

func NewMiniRepo(client *minio.Client) *MiniRepo {
	return &MiniRepo{
		client:     client,
		bucketName: config.Cfg.Minio.DefaultBucket,
	}
}

func (r *MiniRepo) DeleteFile(ctx context.Context, objectName string) error {
	return r.client.RemoveObject(ctx, r.bucketName, objectName, minio.RemoveObjectOptions{})
}

// GetPresignedDownloadURL Get presigned download URL, which can be used to download the file
func (r *MiniRepo) GetPresignedDownloadURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := r.client.PresignedGetObject(ctx, r.bucketName, objectName, expiry, reqParams)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// GetPresignedPostFolderUploadURL Get presigned post folder upload URL, which can be used to upload files to the specified folder
func (r *MiniRepo) GetPresignedPostFolderUploadURL(ctx context.Context, folderPath string, expiry time.Duration) (string, map[string]string, error) {
	policy := minio.NewPostPolicy()
	err := policy.SetBucket(r.bucketName)
	if err != nil {
		return "", nil, fmt.Errorf("failed to set bucket: %w", err)
	}

	err = policy.SetKeyStartsWith(folderPath + "/")
	if err != nil {
		return "", nil, fmt.Errorf("failed to set key prefix: %w", err)
	}

	err = policy.SetExpires(time.Now().Add(expiry))
	if err != nil {
		return "", nil, fmt.Errorf("failed to set expiration: %w", err)
	}

	// policy.SetContentLengthRange(1024, 10*1024*1024)
	postUrl, formData, err := r.client.PresignedPostPolicy(ctx, policy)
	return postUrl.String(), formData, err
}

// DeleteFolder Delete folder
func (r *MiniRepo) DeleteFolder(ctx context.Context, folderPath string) error {
	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	objectsCh := r.client.ListObjects(ctx, r.bucketName, minio.ListObjectsOptions{
		Prefix:    folderPath,
		Recursive: true,
	})

	errorCh := r.client.RemoveObjects(ctx, r.bucketName, objectsCh, minio.RemoveObjectsOptions{})

	for e := range errorCh {
		if e.Err != nil {
			return fmt.Errorf("failed to remove object %s: %w", e.ObjectName, e.Err)
		}
	}
	return nil
}
