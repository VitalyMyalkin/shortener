package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/VitalyMyalkin/shortener/internal/config"
	"github.com/VitalyMyalkin/shortener/internal/storage"
	"github.com/VitalyMyalkin/shortener/internal/logger"
)

type App struct {
	Cfg     config.Config
	Storage *storage.Storage
	short   int
}

type Request struct {
    URLstring string    `json:"url"`
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

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	url, err := url.ParseRequestURI(string(body))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": string(body) + "не является валидным URL",
		})
	}
	newApp.short += 1
	newApp.Storage.AddOrigin(strconv.Itoa(newApp.short), url)

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
		logger.Log.Debug(req.URLstring + "не является валидным URL", zap.Error(err))
		c.String(http.StatusBadRequest, "")
	}
	newApp.short += 1
	newApp.Storage.AddOrigin(strconv.Itoa(newApp.short), url)

	c.Header("content-type", "application/json")
	
	c.JSON(http.StatusCreated, gin.H{
		"result": newApp.Cfg.ShortenAddr+"/"+strconv.Itoa(newApp.short),
	})
}

func (newApp *App) GetOrigin(c *gin.Context) {

	original, ok := newApp.Storage.Storage[c.Param("id")]
	if ok {
		c.Header("Location", original)
		c.Status(http.StatusTemporaryRedirect)
	} else {
		c.Status(http.StatusNotFound)
	}
}
