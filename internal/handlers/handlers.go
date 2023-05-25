package handlers

import (
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VitalyMyalkin/shortener/internal/config"
	"github.com/VitalyMyalkin/shortener/internal/storage"
)

type App struct {
	Cfg     config.Config
	Storage *storage.Storage
	short   int
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

func (newApp *App) GetOrigin(c *gin.Context) {

	original, ok := newApp.Storage.Storage[c.Param("id")]
	if ok {
		c.Header("Location", original)
		c.Status(http.StatusTemporaryRedirect)
	} else {
		c.Status(http.StatusNotFound)
	}
}
