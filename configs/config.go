package configs

import "github.com/nazarkorpal/img-task/tools"

type Config struct {
	DB    *DB
	RABBIT_URL string
}

type DB struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
}

func New() *Config {
	return &Config{
		DB: &DB{
			DB_USER:     tools.GetEnv("DB_USER", ""),
			DB_HOST:     tools.GetEnv("DB_HOST", ""),
			DB_PASSWORD: tools.GetEnv("DB_PASSWORD", ""),
		},
		RABBIT_URL: tools.GetEnv("RABBIT_URL", ""),
	}
}
