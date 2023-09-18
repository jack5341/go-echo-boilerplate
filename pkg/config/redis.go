package config

type Redis struct {
	REDIS_HOST     string `env:"REDIS_HOST"`
	REDIS_PORT     string `env:"REDIS_PORT"`
	REDIS_PASSWORD string `env:"REDIS_PASSWORD"`
}
