package s3storage

import (
	"ChatServer/pkg/repository/minio/internal"
	"context"
	"github.com/minio/minio-go/v7"
	"io"
)

type Object interface {
	PutObject(ctx context.Context, objectName string, reader io.Reader, objSize int64) (string, error)
	GetObject(ctx context.Context, objectName string) (*minio.Object, error)
}

type Bucket interface {
	CreateBucket(ctx context.Context, bucketName string) error
}

type S3Storage struct {
	Object
	Bucket
}

func NewS3Storage(client *minio.Client) *S3Storage {
	return &S3Storage{
		Object: internal.NewS3Object(client),
		Bucket: internal.NewS3Bucket(client),
	}
}
