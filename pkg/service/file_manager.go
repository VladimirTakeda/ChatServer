package service

import (
	"ChatServer/pkg/repository/minio"
	"context"
	"github.com/minio/minio-go/v7"
	"io"
	"strconv"
	"strings"
)

type FileManagerService struct {
	s3Storage *s3storage.S3Storage
}

func NewFileManagerService(s3Storage *s3storage.S3Storage) *FileManagerService {
	return &FileManagerService{s3Storage: s3Storage}
}

func (s *FileManagerService) SaveFile(ctx context.Context, fromId, chatId int, file io.Reader, size int64, fileName string) (string, error) {
	// generate unique object name, for example (userId + chatId + fileName)
	objName := strings.Join([]string{strconv.Itoa(fromId), strconv.Itoa(chatId), fileName}, "")
	_, err := s.s3Storage.PutObject(ctx, objName, file, size)
	if err != nil {
		return "", err
	}
	return objName, nil
}

func (s *FileManagerService) LoadFile(ctx context.Context, objName string) (*minio.Object, error) {
	obj, err := s.s3Storage.GetObject(ctx, objName)
	if err != nil {
		return nil, err
	}
	return obj, nil
}
