package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jantoniogonzalez/factos/internal/handlers"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	var router *chi.Mux = chi.NewRouter()
	handlers.Handlers(router)

	fmt.Printf(("Starting GO API service..."))

	err := http.ListenAndServe(":3333", router)

	if err != nil {
		log.Error(err)
	}
}
