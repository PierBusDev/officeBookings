package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf provides CSRF protection for every request
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appConfig.InProduction, //TODO change when running under https
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saves the session on every req
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
