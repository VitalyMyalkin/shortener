package main

import (
	"io"
	"net/http"
)

func postHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только Post запросы!!", http.StatusBadRequest)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	w.Header().Set("content-type", "text/plain")
	// устанавливаем код 201
	w.WriteHeader(http.StatusCreated)
	// пишем тело ответа
	w.Write([]byte(body))
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только Get запросы!!", http.StatusBadRequest)
		return
	}
	id := r.URL.Query().Get("id")

	w.Header().Set("Location", id)
	// устанавливаем код 307
	w.WriteHeader(http.StatusTemporaryRedirect)
	// пишем тело ответа
	w.Write([]byte(""))
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
