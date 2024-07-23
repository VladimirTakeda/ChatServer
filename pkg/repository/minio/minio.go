package s3storage

import (
	"ChatServer/pkg/repository/minio/internal"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type Config struct {
	EndPoint  string
	KeyID     string
	AccessKey string
}

func NewMinioClient(cfg Config) (*minio.Client, error) {
	// Инициализация клиента MinIO
	client, err := minio.New(cfg.EndPoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.KeyID, cfg.AccessKey, ""),
		Secure: false,
	})
	if err != nil {
		logrus.Fatalln("can't create minio client: " + err.Error())
		return nil, err
	}

	// Создание бакета
	ctx := context.Background()
	exists, errBucketExists := client.BucketExists(ctx, internal.DefaultBucketName)
	if errBucketExists != nil {
		logrus.Debugln("can't find minio bucket: " + errBucketExists.Error())
	}

	if !exists {
		err = client.MakeBucket(ctx, internal.DefaultBucketName, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Fatalln("can't create minio bucket: " + err.Error())
			return nil, err
		}
	}

	return client, nil
}
