package main

import (
	"io"
	"net/http"
)

type MyMap map[string]string

var m MyMap

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только Post запросы!!", http.StatusBadRequest)
		return
	}

	m := make(MyMap)
	var i int

	body, err := io.ReadAll(r.Body)
	i += 1
	if err != nil {
		return
	}

	m[string(rune(i))] = string(body)

	answer := "http://localhost:8080/" + string(rune(i))

	w.Header().Set("content-type", "text/plain")
	// устанавливаем код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte(answer))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только Get запросы!!", http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")

	original := m[id]

	w.Header().Set("Location", original)
	// устанавливаем код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
	// пишем тело ответа
	io.WriteString(w, "")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, postHandler)
	mux.HandleFunc(`/{id}`, getHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
