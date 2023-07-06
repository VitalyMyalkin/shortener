package handlers

import (
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"bufio"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/VitalyMyalkin/shortener/internal/compress"
	"github.com/VitalyMyalkin/shortener/internal/config"
	"github.com/VitalyMyalkin/shortener/internal/logger"
	"github.com/VitalyMyalkin/shortener/internal/storage"
)

type App struct {
	Cfg     config.Config
	Storage *storage.Storage
	short   int
}

type Request struct {
	URLstring string `json:"url"`
}

func NewApp() *App {

	cfg := config.GetConfig()

	storage := storage.NewStorage()

	return &App{
		Cfg:     cfg,
		Storage: storage,
		short:   0,
	}
}

func (newApp *App) GetShortened(c *gin.Context) {

	contentEncoding := c.Request.Header.Get("Content-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")
	if sendsGzip {
		// оборачиваем тело запроса в io.Reader с поддержкой декомпрессии
		cr, err := compress.NewCompressReader(c.Request.Body)
		if err != nil {
			c.Writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		// меняем тело запроса на новое
		c.Request.Body = cr
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		gz, err := gzip.NewReader(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
		body, err = io.ReadAll(gz)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
		}
	}
	url, err := url.ParseRequestURI(string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": string(body) + "не является валидным URL",
		})
	}
	newApp.short += 1
	if newApp.Cfg.FilePath == "" {
		newApp.Storage.AddOrigin(strconv.Itoa(newApp.short), url)
	} else {
		fileName := newApp.Cfg.FilePath
		defer os.Remove(fileName)

		Producer, err := storage.NewProducer(fileName)
		if err != nil {
			logger.Log.Fatal("не создан или не открылся файл записи" + fileName)
		}
		defer Producer.Close()
		if err := Producer.WriteShortenedURL(strconv.Itoa(newApp.short), url); err != nil {
			logger.Log.Fatal("запись не внесена в файл")
		}
	}
	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, newApp.Cfg.ShortenAddr+"/"+strconv.Itoa(newApp.short))
}

func (newApp *App) GetShortenedAPI(c *gin.Context) {

	// десериализуем запрос в структуру модели
	logger.Log.Debug("decoding request")
	var req Request
	dec := json.NewDecoder(c.Request.Body)
	if err := dec.Decode(&req); err != nil {
		logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
		c.String(http.StatusInternalServerError, "")
	}

	url, err := url.ParseRequestURI(req.URLstring)
	if err != nil {
		logger.Log.Debug(req.URLstring+"не является валидным URL", zap.Error(err))
		c.String(http.StatusBadRequest, "")
	}
	newApp.short += 1
	if newApp.Cfg.FilePath == "" {
		newApp.Storage.AddOrigin(strconv.Itoa(newApp.short), url)
	} else {
		fileName := newApp.Cfg.FilePath
		defer os.Remove(fileName)

		Producer, err := storage.NewProducer(fileName)
		if err != nil {
			logger.Log.Fatal("не создан или не открылся файл записи" + fileName)
		}
		defer Producer.Close()
		if err := Producer.WriteShortenedURL(strconv.Itoa(newApp.short), url); err != nil {
			logger.Log.Fatal("запись не внесена в файл")
		}
	}
	c.Header("content-type", "application/json")

	c.JSON(http.StatusCreated, gin.H{
		"result": newApp.Cfg.ShortenAddr + "/" + strconv.Itoa(newApp.short),
	})
}

func (newApp *App) GetOrigin(c *gin.Context) {
	var original string
	
	original = newApp.Storage.Storage[c.Param("id")]
	
	if newApp.Cfg.FilePath != "" {
		fileName := newApp.Cfg.FilePath
		defer os.Remove(fileName)

		file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Log.Fatal("не создан или не открылся файл записи")
		}

		if err != nil {
			logger.Log.Fatal("невозможно прочитать данные файла записи")
		}

		var shortenedURL storage.ShortenedURL

		scanner := bufio.NewScanner(file)
		// optionally, resize scanner's capacity for lines over 64K, see next example
		for scanner.Scan() {
			err = json.Unmarshal(scanner.Bytes(), &shortenedURL)
			if err != nil {
				logger.Log.Fatal("не создана структура")
			} else {
				if shortenedURL.ShortURL == c.Param("id") {
					original = shortenedURL.OriginalURL
				}
			}
		}
	}

	if original != "" {
		c.Header("Location", original)
		c.Status(http.StatusTemporaryRedirect)
	} else {
		c.Status(http.StatusNotFound)
	}
}
