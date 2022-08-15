package storage

import (
	"github.com/nazarkorpal/img-task/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type imageStorage struct {
	db       *gorm.DB
	rabbitMQ *amqp.Channel
}

type Image interface {
	PublishToQueue(queueName string, imageHash string) error
	TakeFromQueue(queueName string) (<-chan amqp.Delivery, error)
	InsertImage(image *models.Image) error
	GetImageByID(imageID uint) (*models.Image, error)
	GetImageByHash(imageHash string) (*models.Image, error)
	UpdateImage(image *models.Image) error
}

func NewImage(db *gorm.DB, rabbitMQ *amqp.Channel) Image {
	return &imageStorage{
		db,
		rabbitMQ,
	}
}

func (i *imageStorage) PublishToQueue(queueName string, imageHash string) error {
	return i.rabbitMQ.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(imageHash),
	})
}

func (i *imageStorage) TakeFromQueue(queueName string) (<-chan amqp.Delivery, error) {
	return i.rabbitMQ.Consume(queueName, "", true, true, true, true, nil)
}

func (i *imageStorage) InsertImage(image *models.Image) error {
	return i.db.Create(image).Error
}

func (i *imageStorage) GetImageByID(imageID uint) (*models.Image, error) {
	var image models.Image
	err := i.db.First(&image, imageID).Error

	return &image, err
}

func (i *imageStorage) GetImageByHash(imageHash string) (*models.Image, error) {
	var image models.Image
	err := i.db.First(&image, "hash = ?", imageHash).Error

	return &image, err
}

func (i *imageStorage) UpdateImage(image *models.Image) error {
	return i.db.Save(image).Error
}
