package main

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

	original, ok := newApp.m[c.Param("id")]
	if ok {
		c.Header("Location", original)
		c.Status(http.StatusTemporaryRedirect)
	} else {
		c.Status(http.StatusNotFound)
	} 
}

func main() {

	newApp := NewApp()

	router := gin.Default()
	router.POST("/", newApp.GetShortened)
	router.GET("/:id", newApp.GetOrigin)

	router.Run(newApp.Cfg.RunAddr)
}
