package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type MinioService struct {
	client     *minio.Client
	bucketName string
}

func NewMinioService() (*MinioService, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	accessKeyID := os.Getenv("MINIO_ACCESS_KEY")
	if accessKeyID == "" {
		accessKeyID = "minioadmin"
	}

	secretAccessKey := os.Getenv("MINIO_SECRET_KEY")
	if secretAccessKey == "" {
		secretAccessKey = "minioadmin"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Minio client: %w", err)
	}

	bucketName := os.Getenv("MINIO_BUCKET_NAME")
	if bucketName == "" {
		bucketName = "acids"
	}

	service := &MinioService{
		client:     client,
		bucketName: bucketName,
	}

	if err := service.ensureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket exists: %w", err)
	}

	return service, nil
}

func (ms *MinioService) ensureBucket() error {
	ctx := context.Background()
	exists, err := ms.client.BucketExists(ctx, ms.bucketName)
	if err != nil {
		return err
	}

	if !exists {
		err = ms.client.MakeBucket(ctx, ms.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return err
		}
		logrus.Infof("Created bucket: %s", ms.bucketName)
	}

	return nil
}

func (ms *MinioService) UploadImage(file *multipart.FileHeader, acidID int) (string, error) {
	ctx := context.Background()

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("acid_%d_%d%s", acidID, time.Now().Unix(), ext)

	_, err = ms.client.PutObject(ctx, ms.bucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	logrus.Infof("Successfully uploaded image: %s", objectName)
	return objectName, nil
}

func (ms *MinioService) DeleteImage(objectName string) error {
	ctx := context.Background()

	if strings.Contains(objectName, "/") {
		parts := strings.Split(objectName, "/")
		objectName = parts[len(parts)-1]
	}

	err := ms.client.RemoveObject(ctx, ms.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	logrus.Infof("Successfully deleted image: %s", objectName)
	return nil
}

func (ms *MinioService) GetImageURL(objectName string) string {
	if strings.Contains(objectName, "/") {
		parts := strings.Split(objectName, "/")
		objectName = parts[len(parts)-1]
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}

	useSSL := os.Getenv("MINIO_USE_SSL") == "true"
	protocol := "http"
	if useSSL {
		protocol = "https"
	}

	return fmt.Sprintf("%s://%s/%s/%s", protocol, endpoint, ms.bucketName, objectName)
}
