package main

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
)

func addSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com; img-src 'self' https://media.api-sports.io; child-src 'self' https://accounts.google.com;")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	csrfHandler.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("Failed CSRF validation")
		w.WriteHeader(http.StatusBadRequest)
	}))

	return csrfHandler
}

// func requireAuthentication() http.Handler {

// }
