package main

import (
	"net/http"

	auth "my-echo-server/services"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello world!")
	})
	g := e.Group("/a")
	auth.Register(g)
	e.Logger.Fatal((e.Start(":3000")))
}
