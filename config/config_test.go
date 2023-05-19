package config

import (
	"os"
	"testing"
	"flag"

	"github.com/stretchr/testify/assert"
)

func TestConfig_parseFlags(t *testing.T) {
	cfg := Config{}
	os.Args = []string{"cmd", "-а", "localhost:8082", "-b", "http://localhost:8081"}
	cfg.parseFlags()
	assert.Equal(t, cfg.RunAddr, "localhost:8082", "Адрес сервера не совпадает с ожидаемым")
	assert.Equal(t, cfg.ShortenAddr, "http://localhost:8081", "Короткая ссылка не совпадает с ожидаемым")
}

func TestGetConfig(t *testing.T) {
	cfg := Config{}
	os.Setenv("SERVER_ADDRESS", "localhost:8082")
	os.Setenv("BASE_URL", "http://localhost:8081")
	GetConfig()
	assert.Equal(t, cfg.RunAddr, "localhost:8082", "Адрес сервера не совпадает с ожидаемым")
	assert.Equal(t, cfg.ShortenAddr, "http://localhost:8081", "Короткая ссылка не совпадает с ожидаемым")
}
