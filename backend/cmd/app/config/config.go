package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on real environment variables")
	}

	if err := env.Parse(&Config.Application); err != nil {
		panic(err)
	}
	if err := env.Parse(&Config.Database); err != nil {
		panic(err)
	}
	if err := env.Parse(&Config.Auth); err != nil {
		panic(err)
	}

}
