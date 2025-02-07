package main

import "github.com/julienschmidt/httprouter"

func newRouter() *httprouter.Router {
	router := httprouter.New()
	return router
}
