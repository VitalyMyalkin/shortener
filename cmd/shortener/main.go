package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/VitalyMyalkin/shortener/internal/compress"
	"github.com/VitalyMyalkin/shortener/internal/handlers"
	"github.com/VitalyMyalkin/shortener/internal/logger"
)

func main() {

	newApp := handlers.NewApp()

	router := gin.Default()
	logger.Initialize()
	router.Use(logger.WithLogging())
	router.Use(compress.GzipMiddleware())

	router.POST("/", newApp.GetShortened)
	router.POST("/api/shorten", newApp.GetShortenedAPI)
	router.GET("/:id", newApp.GetOrigin)
	router.GET("/ping", newApp.PingPostgresDB)
	router.POST("/api/shorten/batch", newApp.SendBatch)

	if err := router.Run(newApp.Cfg.RunAddr); err != nil {
		fmt.Println(err)
	}
}
