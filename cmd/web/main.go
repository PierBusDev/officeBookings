package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pierbusdev/conferenceRoomBookings/internal/config"
	"github.com/pierbusdev/conferenceRoomBookings/internal/handlers"
	"github.com/pierbusdev/conferenceRoomBookings/internal/helpers"
	"github.com/pierbusdev/conferenceRoomBookings/internal/models"
	"github.com/pierbusdev/conferenceRoomBookings/internal/render"
)

const portNumber = ":4554"

var appConfig config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Starting application on port %s\n", portNumber)
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&appConfig),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting the server: ", err)
	}
}

func run() error {
	//this below is needed to be capable of storing complex types in the session
	gob.Register(models.Reservation{})
	//creating app config
	//TODO change this to true when in production
	appConfig.InProduction = false

	//creating logs
	infoLog = log.New(os.Stdout, "INFO =>\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog
	errorLog = log.New(os.Stdout, "ERROR =>\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	//initializing session manager
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //cookie must persist after user closes the browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction //TODO change this in production

	appConfig.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Could not create template cache:", err)
		return err
	}
	//initialize values of appConfig
	appConfig.TemplateCache = templateCache
	appConfig.UseCache = false
	//passing config to the template package
	render.NewTemplate(&appConfig)

	//creating and passing a new Repository to the handlers package
	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

	//helpers
	helpers.NewHelpers(&appConfig)

	return nil
}
