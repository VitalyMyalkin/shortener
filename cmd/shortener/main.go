package main

import (
	"github.com/gin-gonic/gin"

	"github.com/VitalyMyalkin/shortener/internal/handlers"
)

func main() {

	newApp := handlers.NewApp()

	router := gin.Default()
	router.POST("/", newApp.GetShortened)
	router.GET("/:id", newApp.GetOrigin)

	router.Run(newApp.Cfg.RunAddr)
}
