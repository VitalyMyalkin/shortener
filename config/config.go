package config

import (
	"flag"

	"github.com/caarlos0/env/v8"
)

// неэкспортированная переменная runAddr содержит адрес и порт для запуска сервера
type Config struct {
	RunAddr     string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	ShortenAddr string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func (cfg Config) parseFlags() {

	flag.StringVar(&cfg.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&cfg.ShortenAddr, "b", "http://localhost:8080", "default part of shortened URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()
}

func GetConfig() Config {

	cfg := Config{}

	cfg.parseFlags()

	env.Parse(&cfg)

	return cfg
}
