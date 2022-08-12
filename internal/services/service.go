package services

import "github.com/nazarkorpal/img-task/internal/storage"

type Service struct{
	Image
}

func NewService(storage *storage.Storage) *Service {
	return &Service{
		Image: NewImage(storage.Image),
	}
}
