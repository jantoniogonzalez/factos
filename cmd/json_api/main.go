package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	port := ":8080"

	srv := &http.Server{
		Addr:         port,
		Handler:      newRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
