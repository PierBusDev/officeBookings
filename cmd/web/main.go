package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/pierbusdev/basicWeb/pkg/config"
	"github.com/pierbusdev/basicWeb/pkg/handlers"
	"github.com/pierbusdev/basicWeb/pkg/render"
)

const portNumber = ":4554"

func main() {
	//creating app config
	var appConfig config.AppConfig
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
