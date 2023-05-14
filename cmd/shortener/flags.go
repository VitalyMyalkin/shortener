package main

import (
	"flag"
	"os"
)

// неэкспортированная переменная runAddr содержит адрес и порт для запуска сервера
var runAddr string
var shortenAddr string

// parseFlags обрабатывает аргументы командной строки
// и сохраняет их значения в соответствующих переменных
func parseFlags() {
	// регистрируем переменную runAddr
	// как аргумент -a со значением :8080 по умолчанию
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&shortenAddr, "b", "http://localhost:8080", "default part of shortened URL")
	// парсим переданные серверу аргументы в зарегистрированные переменные
	flag.Parse()

	if envRunAddr := os.Getenv("SERVER_ADDRESS"); envRunAddr != "" {
        runAddr = envRunAddr
    }
	if envShortenAddr := os.Getenv("BASE_URL"); envShortenAddr != "" {
        shortenAddr = envShortenAddr
    }
}

