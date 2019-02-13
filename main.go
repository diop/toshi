package main

import (
	"net/http"
	"os"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	port := os.Getenv("PORT")
	if port == "" {
		e.Logger.Fatal("$PORT must be set")
	}

	// Set up Echo, configure server side validation, and hook into middleware.
	e.Server.Addr = ":" + port
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello from Toshi!")
	})

	// Gracefully shut down the server on interrupt.
	e.Logger.Fatal(gracehttp.Serve(e.Server))
}
