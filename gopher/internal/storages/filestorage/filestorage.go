package filestorage

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	*minio.Client
	bucketname string
}

func (s *FileStorage) SaveFile(ctx context.Context, objectName string, r io.Reader, size int64) (string, error) {
	const op = "minio.SaveFile"

	_, err := s.PutObject(
		ctx,
		s.bucketname,
		objectName,
		r,    // <- твой io.Reader
		size, // <- важно!
		minio.PutObjectOptions{
			ContentType: "audio/webm",
		},
	)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return s.bucketname + objectName, nil
}

func NewFileStorage(address, username, password, bucketname string) (*FileStorage, error) {
	const op = "minio.NewFileStorage"

	client, err := minio.New(address, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	exists, err := client.BucketExists(context.Background(), bucketname)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !exists {
		return nil, fmt.Errorf("%s: %w", op, fmt.Errorf("bucket %s does not exist", bucketname))
	}

	return &FileStorage{client, bucketname}, nil
}
