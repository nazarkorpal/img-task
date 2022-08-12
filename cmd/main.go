package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nazarkorpal/img-task/internal/api"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {
	api.New().Start()
}
