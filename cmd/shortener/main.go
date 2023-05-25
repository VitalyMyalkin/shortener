package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/VitalyMyalkin/shortener/internal/handlers"
)

func main() {

	newApp := handlers.NewApp()

	router := gin.Default()
	router.POST("/", newApp.GetShortened)
	router.GET("/:id", newApp.GetOrigin)

	if err := router.Run(newApp.Cfg.RunAddr); err != nil {
		fmt.Println(err)
	}
}
