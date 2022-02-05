package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/serjyuriev/shortener/internal/app"
)

func main() {
	c := app.Config{
		UrlLength: 4,
		Host:      "localhost",
		Port:      8080,
	}
	svc := app.MakeService(c)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", c.Host, c.Port), svc))
}
