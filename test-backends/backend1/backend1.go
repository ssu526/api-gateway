package main

import (
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("backend1", "true")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from backend 1\n"))
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/", hello)

	server := &http.Server{
		Addr:    ":9091",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting the server for backend 1")
	}
}
