package setup

import (
	"fmt"

	"github.com/nazarkorpal/img-task/configs"
	"github.com/nazarkorpal/img-task/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(conf *configs.DB) (*gorm.DB, error) {
	var db *gorm.DB

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres sslmode=disable", conf.DB_HOST, conf.DB_USER, conf.DB_PASSWORD)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&models.Image{})

	return db, err
}
