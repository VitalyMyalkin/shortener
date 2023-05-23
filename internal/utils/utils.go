package utils

import (
	"net/url"
	"strconv"
)

type Storage struct {
	Storage map[string]string
}

func NewStorage() *Storage {

	storage := make(map[string]string)

	return &Storage{
		Storage: storage,
	}
}

func (newStorage *Storage) AddOrigin(short int, url *url.URL) {

	newStorage.Storage[strconv.Itoa(short)] = url.String()
}
