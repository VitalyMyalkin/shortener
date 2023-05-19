package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/VitalyMyalkin/shortener/internal/handlers"
)

func Test_getShortened(t *testing.T) {
	origin := "https://practicum.yandex.ru/"
	r := gin.Default()
	newApp := handlers.NewApp()
	r.POST("/", newApp.GetShortened)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(origin)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Код ответа не совпадает с ожидаемым")
	// assert.Equal(t, origin, newApp.m["1"], "В базе не появилась запись")
	assert.Equal(t, newApp.Cfg.ShortenAddr+"/1", w.Body.String(), "Тело ответа не совпадает с ожидаемым")
}

func Test_getOrigin(t *testing.T) {
	r := gin.Default()
	newApp := handlers.NewApp()
	// newApp.m["1"] = "https://practicum.yandex.ru/"
	r.GET("/1", newApp.GetOrigin)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/1", bytes.NewBuffer([]byte("")))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusTemporaryRedirect, w.Code, "Код ответа не совпадает с ожидаемым")
	assert.Equal(t, "", w.Body.String(), "Тело ответа не совпадает с ожидаемым")
}
