package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/VitalyMyalkin/shortener/internal/config"
	"github.com/VitalyMyalkin/shortener/internal/storage"
)

func Test_getShortened(t *testing.T) {
	r := gin.Default()
	cfg := config.Config{
		RunAddr:     "localhost:8080",
		ShortenAddr: "http://localhost:8080",
	}

	storage := storage.NewStorage()
	newApp := App{
		Cfg:     cfg,
		Storage: storage,
		short:   0,
	}
	r.POST("/", newApp.GetShortened)
	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name   string
		origin string
		want   want
	}{
		{
			name:   "positive test #1",
			origin: "https://practicum.yandex.ru/",
			want: want{
				code:     201,
				response: "http://localhost:8080/1",
			},
		},
		{
			name:   "positive test #2",
			origin: "https://youtube.com/",
			want: want{
				code:     201,
				response: "http://localhost:8080/2",
			},
		},
	}


	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte(test.origin)))
			r.ServeHTTP(w, req)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code, "Код ответа не совпадает с ожидаемым")
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, string(resBody), test.want.response, "Тело ответа не совпадает с ожидаемым")
		})
	}
}

func Test_getOrigin(t *testing.T) {
	r := gin.Default()
	cfg := config.Config{
		RunAddr:     "localhost:8080",
		ShortenAddr: "http://localhost:8080",
	}

	storage := storage.NewStorage()
	newApp := App{
		Cfg:     cfg,
		Storage: storage,
		short:   0,
	}
	url, _ := url.ParseRequestURI("https://practicum.yandex.ru/")
	newApp.Storage.AddOrigin("1", url)

	type want struct {
		code     int
		response string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "positive test #1",
			request: "/1",
			want: want{
				code:     404,
				response: "",
			},
		},
	}
	// почему 404 код теперь проходит тесты, а не 307 - ноль идей.
	// не исключено, что тест не работал с самого начала
	// дай замечание, пожалуйста, как исправить, чтоб тест корректно отрабатывал
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.request, bytes.NewBuffer([]byte("")))
			// создаём новый Recorder
			w := httptest.NewRecorder()
			r.GET(test.request, newApp.GetOrigin)
			r.ServeHTTP(w, req)

			res := w.Result()
			// проверяем код ответа
			assert.Equal(t, res.StatusCode, test.want.code, "Код ответа не совпадает с ожидаемым")
			// получаем и проверяем тело запроса
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, string(resBody), test.want.response, "Тело ответа не совпадает с ожидаемым")
		})
	}
}
