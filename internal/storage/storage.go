package storage

import (
	"os"
	"encoding/json"
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

type ShortenedURL struct {
    ID       string    `json:"uuid"`
    ShortURL string  `json:"short_url"`
    OriginalURL    string `json:"original_url"`
}

type FileWriter struct {
    file    *os.File
}

func NewFileWriter(fileName string) (*FileWriter, error) {
    file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
    if err != nil {
        return nil, err
    }

    return &FileWriter{
        file:    file,
    }, nil
}

func (p *FileWriter) WriteShortenedURL(short string, url *url.URL) error {
    
	shortenedURL := ShortenedURL{
		ID:     short,
		ShortURL: short,
		OriginalURL :      url.String(), 
	}
    data, err := json.Marshal(&shortenedURL)
    if err != nil {
        return err
    }
    // добавляем перенос строки
    data = append(data, '\n')

    _, err = p.file.Write(data)
    return err
}

func (p *FileWriter) Close() error {
    return p.file.Close()
}




