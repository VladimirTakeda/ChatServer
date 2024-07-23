package internal

import (
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
)

type S3Object struct {
	S3Client *minio.Client
}

const DefaultBucketName = "files"

func NewS3Object(s3Client *minio.Client) *S3Object {
	return &S3Object{S3Client: s3Client}
}

func (r *S3Object) PutObject(ctx context.Context, objectName string, reader io.Reader, objSize int64) (string, error) {
	_, err := r.S3Client.PutObject(ctx, DefaultBucketName, objectName, reader, objSize, minio.PutObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	return objectName, nil
}

func (r *S3Object) GetObject(ctx context.Context, objectName string) (*minio.Object, error) {
	Object, err := r.S3Client.GetObject(ctx, DefaultBucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	return Object, nil
}
