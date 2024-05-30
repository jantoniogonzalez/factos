package handlers

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func Handlers(router chi.Router) {
	router.Use(chiMiddleware.StripSlashes)

	// router.Route("/factos", func(router chi.Router) {
	// 	router.Get("/")     // GET all factos
	// 	router.Post("/new") // CREATE a new factos

	// 	router.Route("/{factoId}", func(router chi.Router) {
	// 		router.Get("/")    // GET facto with id
	// 		router.Put("/")    // EDIT an existing facto
	// 		router.Delete("/") // DELETE an existing facto

	// 	})
	// })
}
