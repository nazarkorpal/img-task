package storage

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

type Storage struct {
	Image
}

func NewStorage(db *gorm.DB, rabbitMQ *amqp.Channel) *Storage {
	return &Storage{
		Image: NewImage(db, rabbitMQ),
	}
}
