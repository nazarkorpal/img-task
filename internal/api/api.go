package api

import (
	"log"

	"github.com/nazarkorpal/img-task/configs"
	"github.com/nazarkorpal/img-task/internal/handlers"
	"github.com/nazarkorpal/img-task/internal/services"
	"github.com/nazarkorpal/img-task/internal/storage"
	"github.com/nazarkorpal/img-task/setup"
)

type API struct {
	Config *configs.Config
}

func New() *API {
	return &API{
		Config: configs.New(),
	}
}

func (a *API) Start() {
	logger := log.Default()
	db, err := setup.ConnectDB(a.Config.DB)
	if err != nil {
		log.Panic(err)
	}
	logger.Println("DB connected successfully")

	rabbitMQ, err := setup.ConnectRabbitMQ(a.Config.RABBIT_URL)
	if err != nil {
		log.Panic(err)
	}
	logger.Println("RabbitMQ connected successfully")

	storage := storage.NewStorage(db, rabbitMQ)
	services := services.NewService(storage)
	handlers := handlers.NewHandler(services)

	handlers.Init().Run()
}
