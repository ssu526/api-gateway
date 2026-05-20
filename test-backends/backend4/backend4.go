package main

import (
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("backend4", "true")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from backend 4\n"))
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
		Addr:    ":9094",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting the server for backend 4")
	}
}
