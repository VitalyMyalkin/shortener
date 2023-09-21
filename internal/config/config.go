package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	RunAddr     	string `env:"SERVER_ADDRESS"`
	ShortenAddr 	string `env:"BASE_URL"`
	FilePath 		string `env:"FILE_STORAGE_PATH"`
	PostgresDBAddr	string `env:"DATABASE_DSN"`
}

func GetConfig() Config {

	cfg := Config{}

	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.ShortenAddr, "b", "http://localhost:8080", "default part of shortened URL")
	flag.StringVar(&cfg.FilePath, "f", "/tmp/short-url-db.json", "default path to save shortened URL")
	flag.StringVar(&cfg.PostgresDBAddr, "d", "", "default database path")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	env.Parse(&cfg)

	return cfg
}
