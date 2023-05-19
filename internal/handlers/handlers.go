package handlers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VitalyMyalkin/shortener/config"
)

type App struct {
	Cfg config.Config
	m   map[string]string
	i   int
}

func NewApp() *App {

	cfg := config.GetConfig()

	m := make(map[string]string)

	return &App{
		Cfg: cfg,
		m:   m,
		i:   0,
	}
}

func (newApp App) GetShortened(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}
	newApp.i += 1
	newApp.m[strconv.Itoa(newApp.i)] = string(body)

	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, newApp.Cfg.ShortenAddr+"/"+strconv.Itoa(newApp.i))
}

func (newApp App) GetOrigin(c *gin.Context) {

	original := newApp.m[c.Param("id")]
	if original == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Длинной ссылки по заданной короткой ссылке не существует!",
		})
	}

	c.Header("Location", original)
	c.Status(http.StatusTemporaryRedirect)
}
