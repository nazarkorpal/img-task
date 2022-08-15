package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/nazarkorpal/img-task/internal/api"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}
}

func main() {
	api.New().Start()
}
