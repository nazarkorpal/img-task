package storage

import (
	"github.com/nazarkorpal/img-task/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type imageStorage struct {
	db *gorm.DB
	rabbitMQ *amqp.Channel
}

type Image interface {
	PublishToQueue(queueName string, imageHash string) error
	TakeFromQueue(queueName string) (string, error)
	InsertImage(image *models.Image) error
}

func NewImage(db *gorm.DB, rabbitMQ *amqp.Channel) Image {
	return &imageStorage{
		db,
		rabbitMQ,
	}
}

func(i *imageStorage) PublishToQueue(queueName string, imageHash string) error {
	return i.rabbitMQ.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body: []byte(imageHash),
	})
}

func(i *imageStorage) TakeFromQueue(queueName string) (string, error) {
	ch, err := i.rabbitMQ.Consume(queueName, "", true, false, false, true, nil)
	msg := <-ch
	return string(msg.Body), err
}

func(i *imageStorage) InsertImage(image *models.Image) error {
	return i.db.Create(image).Error
}