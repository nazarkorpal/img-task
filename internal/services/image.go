package services

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/nazarkorpal/img-task/internal/models"
	"github.com/nazarkorpal/img-task/internal/storage"
	"github.com/nazarkorpal/img-task/tools"
)

const (
	QUEUE_NAME = "images"
)

type imageService struct {
	ImageStorage storage.Image
}

type Image interface {
	AddImage(file multipart.File, fileHeader *multipart.FileHeader) error
}

func NewImage(imageStorage storage.Image) Image {
	return &imageService{ImageStorage: imageStorage}
}

func (i *imageService) AddImage(file multipart.File, fileHeader *multipart.FileHeader) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	hash := uuid.New().String()
	ext := path.Ext(fileHeader.Filename)

	path := path.Join(dir, "uploads", hash+ext)

	bytes, err := tools.ReadFile(file)
	if err != nil {
		return fmt.Errorf("cannot read file %w", err)
	}

	if err := tools.WriteFile(path, bytes); err != nil {
		return fmt.Errorf("cannot write file %w", err)
	}

	image := &models.Image{
		Hash: hash,
		Ext:  ext,
		Name: strings.Split(fileHeader.Filename, ".")[0],
	}

	if err := i.ImageStorage.InsertImage(image); err != nil {
		return err
	}

	return i.ImageStorage.PublishToQueue(QUEUE_NAME, hash)
}
