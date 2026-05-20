package main

import (
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("backend2", "true")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from backend 2\n"))
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", hello)
	router.HandleFunc("/health", health)

	server := &http.Server{
		Addr:    ":9092",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting the server for backend 2")
	}
}
