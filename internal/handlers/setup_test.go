package handlers

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/justinas/nosurf"
	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
	"github.com/pierbusdev/conferenceRoomBookings/internal/render"
)

var appConfig config.AppConfig
var session *scs.SessionManager

func getRoutes() http.Handler {
	//this below is needed to be capable of storing complex types in the session
	gob.Register(models.Reservation{})
	//creating app config
	//TODO change this to true when in production
	appConfig.InProduction = false

	//creating logs
	infoLog := log.New(os.Stdout, "INFO =>\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog
	errorLog := log.New(os.Stdout, "ERROR =>\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //cookie must persist after user closes the browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction //TODO change this in production

	appConfig.Session = session

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("Could not create template cache:", err)
	}
	//initialize values of appConfig
	appConfig.TemplateCache = templateCache
	appConfig.UseCache = true
	//passing config to the template package
	render.NewTemplate(&appConfig)

	//creating and passing a new Repository to the handlers package
	repo := NewRepo(&appConfig)
	NewHandlers(repo)

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf) //no need to test here too, already tested in middlewares
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/big-office", Repo.BigOffice)
	mux.Get("/medium-office", Repo.MediumOffice)
	mux.Get("/shared-office", Repo.SharedOffice)

	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJson)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)

	mux.Get("/contact", Repo.Contact)

	mux.Get("/reservation-summary", Repo.ReservationSummary)

	//static files loading
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

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

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

var pathToTemplates = "./../../templates"

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	templateCache := make(map[string]*template.Template)

	//find all the files with .page pseudoextension
	pages, err := filepath.Glob(pathToTemplates + "/*page.html")
	if err != nil {
		return templateCache, err
	}

	//we need to check if there are layouts to add to the various pages
	for _, page := range pages {
		fileName := filepath.Base(page) //basically the filename stripped of the path
		templ, err := template.New(fileName).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		//checking for layouts presence
		matches, err := filepath.Glob(pathToTemplates + "/*.layout.html")
		if err != nil {
			return templateCache, err
		}

		if len(matches) > 0 { //if we have layouts we have to incorporate them in the template
			templ, err = templ.ParseGlob(pathToTemplates + "/*.layout.html")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[fileName] = templ
	}

	return templateCache, nil
}
