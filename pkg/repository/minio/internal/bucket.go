package internal

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
)

type S3Bucket struct {
	S3Client *minio.Client
}

func NewS3Bucket(s3Client *minio.Client) *S3Bucket {
	return &S3Bucket{S3Client: s3Client}
}

func (r *S3Bucket) CreateBucket(ctx context.Context, bucketName string) error {
	err := r.S3Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := r.S3Client.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			fmt.Printf("Bucket '%s' already exists\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		fmt.Printf("Successfully created '%s' bucket\n", bucketName)
	}
	return nil
}
