package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	*minio.Client
	bucketname string
}

func NewMinio(address, username, password, bucketname string) (*Minio, error) {
	const op = "minio.NewMinio"

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

	return &Minio{client, bucketname}, nil
}

func (m *Minio) SaveAudioFile(ctx context.Context, objectName string, r io.Reader, size int64) error {
	const op = "minio.SaveAudioFile"

	_, err := m.PutObject(
		ctx,
		m.bucketname,
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

func (m *Minio) AudioFile(ctx context.Context, objectName string) (io.ReadCloser, error) {
	const op = "minio.AudioFile"

	object, err := m.GetObject(
		ctx,
		m.bucketname,
		objectName,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = object.Stat()
	if err != nil {
		return nil, fmt.Errorf("%s (stat): %w", op, err)
	}

	return object, nil
}
