package main

import (
	"fmt"
	"net/http"

	"github.com/pierbusdev/basicWeb/pkg/handlers"
)

const portNumber = ":4554"

func main() {
	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/About", handlers.About)

	fmt.Printf("Starting application on port %s", portNumber)
	http.ListenAndServe(portNumber, nil)
}
