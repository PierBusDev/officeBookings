package main

import (
	"testing"

	"github.com/go-chi/chi"
	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
)

func TestRoutes(t *testing.T) {
	var appConfig config.AppConfig
	mux := routes(&appConfig)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing, it's the correct one
	default:
		t.Errorf("type is not http.Handler but %T", v)
	}
}
