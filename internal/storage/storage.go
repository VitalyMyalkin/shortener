package storage

import (
	"net/url"
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

func (newStorage *Storage) AddOrigin(short string, url *url.URL) {

	newStorage.Storage[short] = url.String()
}
