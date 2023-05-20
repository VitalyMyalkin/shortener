package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/VitalyMyalkin/shortener/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_getShortened(t *testing.T) {
	origin := "https://practicum.yandex.ru/"
	r := gin.Default()
	cfg := config.Config{
		RunAddr:     "localhost:8080",
		ShortenAddr: "http://localhost:8080",
	}

	m := make(map[string]string)
	newApp := App{
		Cfg: cfg,
		m:   m,
		i:   0,
	}
	r.POST("/", newApp.GetShortened)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(origin)))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code, "Код ответа не совпадает с ожидаемым")
	assert.Equal(t, origin, newApp.m["1"], "В базе не появилась запись")
	assert.Equal(t, newApp.Cfg.ShortenAddr+"/1", w.Body.String(), "Тело ответа не совпадает с ожидаемым")
}

func Test_getOrigin(t *testing.T) {
	r := gin.Default()
	cfg := config.Config{
		RunAddr:     "localhost:8080",
		ShortenAddr: "http://localhost:8080",
	}

	m := make(map[string]string)
	newApp := App{
		Cfg: cfg,
		m:   m,
		i:   0,
	}
	newApp.m["1"] = "https://practicum.yandex.ru/"
	r.GET("/1", newApp.GetOrigin)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/1", bytes.NewBuffer([]byte("")))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code, "Код ответа не совпадает с ожидаемым")
	assert.Equal(t, "", w.Body.String(), "Тело ответа не совпадает с ожидаемым")
}
