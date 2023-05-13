package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type MyMap map[string]string

var m MyMap

var i string

func getShortened(c *gin.Context) {

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return
	}
	i += "a"
	m[i] = string(body)

	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, defaultShortenAddr+i)
}

func getOrigin(c *gin.Context) {
	original := m[c.Param("id")]

	c.Header("Location", original)
	c.Status(http.StatusTemporaryRedirect)
}

func main() {
	parseFlags()
	m = make(MyMap)

	router := gin.Default()
	router.POST("/", getShortened)
	router.GET("/:id", getOrigin)

	router.Run(flagRunAddr)
}
