package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/h2non/bimg"
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
	AddImage(file multipart.File, fileHeader *multipart.FileHeader) (uint, error)
	GetImage(imageID uint, quality string) (string, string, error)
	GenerateLessQualityImages() error
}

func NewImage(imageStorage storage.Image) Image {
	return &imageService{ImageStorage: imageStorage}
}

func (i *imageService) AddImage(file multipart.File, fileHeader *multipart.FileHeader) (uint, error) {
	dir, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	hash := uuid.New().String()
	ext := path.Ext(fileHeader.Filename)

	path := path.Join(dir, "uploads", hash+ext)

	bytes, err := tools.ReadFile(file)
	if err != nil {
		return 0, fmt.Errorf("cannot read file %w", err)
	}

	if err := tools.WriteFile(path, bytes); err != nil {
		return 0, fmt.Errorf("cannot write file %w", err)
	}

	image := &models.Image{
		Hash: hash,
		Ext:  ext,
		Name: strings.Split(fileHeader.Filename, ".")[0],
	}

	if err := i.ImageStorage.InsertImage(image); err != nil {
		return 0, err
	}

	err = i.ImageStorage.PublishToQueue(QUEUE_NAME, hash)

	return image.ID, err
}

func (i *imageService) GetImage(imageID uint, quality string) (string, string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", "", err
	}

	image, err := i.ImageStorage.GetImageByID(imageID)
	if err != nil {
		return "", "", err
	}

	var hash string
	switch quality {
	case "25":
		hash = image.Hash25
	case "50":
		hash = image.Hash50
	case "75":
		hash = image.Hash75
	case "100":
		hash = image.Hash
	}

	if hash == "" {
		return "", "", errors.New("image in this quality do not exist")
	}

	imageName := image.Name + "_" + quality + image.Ext

	return imageName, path.Join(dir, "uploads", hash+image.Ext), nil
}

func (i *imageService) GenerateLessQualityImages() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	ch, err := i.ImageStorage.TakeFromQueue(QUEUE_NAME)
	if err != nil {
		return err
	}

	generatePath := func(fileName string) string {
		return path.Join(dir, "uploads", fileName)
	}

	for {
		select {
		case data := <-ch:
			image, err := i.ImageStorage.GetImageByHash(string(data.Body))
			if err != nil {
				return err
			}

			hash25 := uuid.New().String()
			hash50 := uuid.New().String()
			hash75 := uuid.New().String()

			image.Hash25 = hash25
			image.Hash50 = hash50
			image.Hash75 = hash75

			bytes, err := os.ReadFile(generatePath(image.Hash + image.Ext))
			if err != nil {
				return err
			}

			//Generating 25% quality image
			newImage, err := bimg.NewImage(bytes).Process(bimg.Options{Quality: 25})
			if err != nil {
				return err
			}

			if err := bimg.Write(generatePath(hash25+image.Ext), newImage); err != nil {
				return err
			}

			//Generating 50% quality image
			newImage, err = bimg.NewImage(bytes).Process(bimg.Options{Quality: 50})
			if err != nil {
				return err
			}

			if err := bimg.Write(generatePath(hash50+image.Ext), newImage); err != nil {
				return err
			}

			//Generating 75% quality image
			newImage, err = bimg.NewImage(bytes).Process(bimg.Options{Quality: 75})
			if err != nil {
				return err
			}

			if err := bimg.Write(generatePath(hash75+image.Ext), newImage); err != nil {
				return err
			}

			//Updating db record
			if err := i.ImageStorage.UpdateImage(image); err != nil {
				return err
			}
		default:
			time.Sleep(time.Second)
		}
	}
}
