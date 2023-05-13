package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type MyMap map[string]string

var m MyMap

var i string

func getShortened(c *gin.Context) {
	type Param struct {
		A string
	}
	param := new(Param)
	i += "a"
	c.Bind(param)
	m[i] = param.A

	c.Header("content-type", "text/plain")
	c.String(http.StatusCreated, "http://localhost:8080/"+i)
}

func getOrigin(c *gin.Context) {
	original := m[c.Param("id")]

	c.Header("Location", original)
	c.Status(http.StatusTemporaryRedirect)
}

func main() {
	m = make(MyMap)

	router := gin.Default()
	router.POST("/", getShortened)
	router.GET("/:id", getOrigin)

	router.Run("localhost:8080")
}
