package config

import (
	"context"
	"log"

	_ "github.com/joho/godotenv/autoload"
	"github.com/sethvargo/go-envconfig"
)

type Configuration struct {
	APP     App
	DB      DB
	AWS     AWS
	Redis   Redis
	DevMode bool
}

func InitConfig() Configuration {
	ctx := context.Background()

	var c Configuration
	if err := envconfig.Process(ctx, &c); err != nil {
		log.Fatal(err)
	}

	return c
}
