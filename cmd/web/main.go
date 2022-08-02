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
	appConfig.TemplateCache = templateCache
	//passing config to the template package
	render.NewTemplate(&appConfig)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/About", handlers.About)

	fmt.Printf("Starting application on port %s", portNumber)
	http.ListenAndServe(portNumber, nil)
}
