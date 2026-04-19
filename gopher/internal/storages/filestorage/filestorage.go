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

func (s *FileStorage) SaveFile(ctx context.Context, objectName string, r io.Reader, size int64) error {
	const op = "minio.SaveFile"

	_, err := s.PutObject(
		ctx,
		s.bucketname,
		objectName,
		r,
		size,
		minio.PutObjectOptions{
			ContentType: "audio/webm",
		},
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *FileStorage) GetFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	const op = "minio.GetFile"

	// Получаем объект из бакета
	object, err := s.GetObject(
		ctx,
		s.bucketname,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Важно: на этом этапе ошибка может не возникнуть,
	// она может проявиться при первом чтении или через Stat()
	_, err = object.Stat()
	if err != nil {
		return nil, fmt.Errorf("%s (stat): %w", op, err)
	}

	return object, nil
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
