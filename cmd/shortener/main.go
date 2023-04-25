package main

import (
	"io"
	"net/http"
	"strings"
)

type MyMap map[string]string

var m MyMap

var i string

func myHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		body, err := io.ReadAll(r.Body)
		i += "a"
		if err != nil {
			return
		}

		m[i] = string(body)

		answer := "http://localhost:8080/" + i

		w.Header().Set("content-type", "text/plain")
		// устанавливаем код 201
		w.WriteHeader(http.StatusCreated)
		// пишем тело ответа
		w.Write([]byte(answer))
	}

	if r.Method == http.MethodGet {

		id := strings.Split(r.URL.Path, "/")

		original := m[id[len(id)-1]]

		w.Header().Set("Location", original)
		// устанавливаем код 307
		w.WriteHeader(http.StatusTemporaryRedirect)
		// пишем тело ответа
		io.WriteString(w, "")
	}

	if (r.Method != http.MethodGet) && (r.Method != http.MethodPost) {
		http.Error(w, "Только Post или Get запросы!!", http.StatusBadRequest)
		return
	}
}

func main() {
	m = make(MyMap)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, myHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
