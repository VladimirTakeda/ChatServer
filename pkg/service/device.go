package service

import (
	"ChatServer/pkg/repository/postgres"
	"context"
	"time"
)

type DeviceService struct {
	repo postgres.Device
}

func NewDeviceService(repo postgres.Device) *DeviceService {
	return &DeviceService{repo: repo}
}

func (s *DeviceService) RegisterDevice(ctx context.Context, deviceId string, userId int) error {
	return s.repo.RegisterDevice(ctx, deviceId, userId)
}

func (s *DeviceService) SaveLastActiveTime(ctx context.Context, userId int, deviceId string, lastTime time.Time) error {
	return s.repo.SaveLastActiveTime(ctx, userId, deviceId, lastTime)
}

func (s *DeviceService) GetLastActiveTime(ctx context.Context, deviceId string, userId int) (time.Time, error) {
	return s.repo.GetLastActiveTime(ctx, deviceId, userId)
}
