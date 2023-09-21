package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/uuid"
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

	short1 := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://practicum.yandex.ru/")).String()
	short2 := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://youtube.com/")).String()

	storage := storage.NewStorage()
	newApp := App{
		Cfg:     cfg,
		Storage: storage,
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
				response: "http://localhost:8080/"+short1,
			},
		},
		{
			name:   "positive test #2",
			origin: "https://youtube.com/",
			want: want{
				code:     201,
				response: "http://localhost:8080/"+short2,
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

func Test_getShortenedAPI(t *testing.T) {
	r := gin.Default()
	cfg := config.Config{
		RunAddr:     "localhost:8080",
		ShortenAddr: "http://localhost:8080",
	}

	short1 := `{"result":"http://localhost:8080/`+uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://practicum.yandex.ru")).String()+`"}`
	short2 := `{"result":"http://localhost:8080/`+uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://youtube.com/")).String()+`"}`

	storage := storage.NewStorage()
	newApp := App{
		Cfg:     cfg,
		Storage: storage,
	}
	r.POST("/api/shorten", newApp.GetShortenedAPI)
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
			origin: `{"url": "https://practicum.yandex.ru"}`,
			want: want{
				code:     201,
				response: short1,
			},
		},
		{
			name:   "positive test #2",
			origin: `{"url": "https://youtube.com/"}`,
			want: want{
				code:     201,
				response: short2,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// создаём новый Recorder
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodPost, "/api/shorten", bytes.NewBuffer([]byte(test.origin)))
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

	short1 := uuid.NewSHA1(uuid.NameSpaceURL, []byte("https://practicum.yandex.ru/")).String()

	storage := storage.NewStorage()
	newApp := App{
		Cfg:     cfg,
		Storage: storage,
	}
	url, _ := url.ParseRequestURI("https://practicum.yandex.ru/")
	newApp.Storage.AddOrigin(short1, url)

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
			request: "/"+short1,
			want: want{
				code:     307,
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
