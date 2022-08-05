package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/pierbusdev/conferenceRoomBookings/pkg/config"
	"github.com/pierbusdev/conferenceRoomBookings/pkg/handlers"
)

func routes(appConfig *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/big-office", handlers.Repo.BigOffice)
	mux.Get("/medium-office", handlers.Repo.MediumOffice)
	mux.Get("/shared-office", handlers.Repo.SharedOffice)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Get("/contact", handlers.Repo.Contact)

	//static files loading
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
