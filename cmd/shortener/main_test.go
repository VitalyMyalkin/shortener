package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MyHandler(t *testing.T) {

	testCases := []struct {
		method       string
		expectedCode int
		expectedBody string
	}{
		{method: http.MethodPost, expectedCode: http.StatusCreated, expectedBody: "http://localhost:8080/a"},
		{method: http.MethodGet, expectedCode: http.StatusTemporaryRedirect, expectedBody: ""},
		{method: http.MethodPut, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{method: http.MethodDelete, expectedCode: http.StatusBadRequest, expectedBody: ""},
	}
	for _, tc := range testCases {
		t.Run(tc.method, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, "/", nil)
			w := httptest.NewRecorder()

			// вызовем хендлер как обычную функцию, без запуска самого сервера
			MyHandler(w, r)

			assert.Equal(t, tc.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			// проверим корректность полученного тела ответа, если мы его ожидаем
			if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
