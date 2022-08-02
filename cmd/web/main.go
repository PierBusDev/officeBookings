package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/pierbusdev/conferenceRoomBookings/pkg/config"
	"github.com/pierbusdev/conferenceRoomBookings/pkg/handlers"
	"github.com/pierbusdev/conferenceRoomBookings/pkg/render"
)

const portNumber = ":4554"

var appConfig config.AppConfig
var session *scs.SessionManager

func main() {
	//creating app config
	//TODO change this to true when in production
	appConfig.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true //cookie must persist after user closes the browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appConfig.InProduction //TODO change this in production

	appConfig.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Could not create template cache:", err)
	}
	//initialize values of appConfig
	appConfig.TemplateCache = templateCache
	appConfig.UseCache = false
	//passing config to the template package
	render.NewTemplate(&appConfig)

	//creating and passing a new Repository to the handlers package
	repo := handlers.NewRepo(&appConfig)
	handlers.NewHandlers(repo)

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
