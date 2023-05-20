package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	RunAddr     string `env:"SERVER_ADDRESS"`
	ShortenAddr string `env:"BASE_URL"`
}

func GetConfig() Config {

	Cfg := Config{}

	flag.StringVar(&Cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Cfg.ShortenAddr, "b", "http://localhost:8080", "default part of shortened URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	env.Parse(&Cfg)

	return Cfg
}
