package handlers

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func Handlers(router chi.Router) {
	router.Use(chiMiddleware.StripSlashes)

	router.Route("/factos", func(router chi.Router) {
		router.Get("/", GetAllFactos)       // GET all factos
		router.Post("/new", CreateNewFacto) // CREATE a new factos

		router.Route("/{factoId}", func(router chi.Router) {
			router.Get("/", GetFactoById)   // GET facto with id
			router.Put("/", EditFacto)      // EDIT an existing facto
			router.Delete("/", DeleteFacto) // DELETE an existing facto

		})

		// router.Route("/{userId}", func (router chi.Router) {
		// 	router.Get("/", GetFactosByUserId)
		// 	router.Get("/correct", GetCorrectFactosByUserId)
		// })
	})
}
