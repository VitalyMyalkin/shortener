package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_getShortened(t *testing.T) {
	origin := "https://practicum.yandex.ru/"
	r := gin.Default()
	m = make(MyMap)
	r.POST("/", getShortened)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(origin)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Код ответа не совпадает с ожидаемым")
	assert.Equal(t, origin, m["a"], "В базе не появилась запись")
	assert.Equal(t, shortenAddr+"/a", w.Body.String(), "Тело ответа не совпадает с ожидаемым")
}

func Test_getOrigin(t *testing.T) {
	r := gin.Default()
	m = make(MyMap)
	m["a"] = "https://practicum.yandex.ru/"
	r.GET("/a", getOrigin)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/a", bytes.NewBuffer([]byte("")))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "Код ответа не совпадает с ожидаемым")
	assert.Equal(t, "", w.Body.String(), "Тело ответа не совпадает с ожидаемым")
}
