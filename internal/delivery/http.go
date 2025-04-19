package delivery

import (
	"net/http"

	_ "github.com/arvaliullin/wapa/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
)

type EchoHttpService struct {
	Echo *echo.Echo
}

func (e *EchoHttpService) Start(address string) {
	e.Echo.Logger.Fatal(e.Echo.Start(address))
}

func NewEchoHttpService() *EchoHttpService {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete},
	}))

	e.GET("/swagger/*", swagger.WrapHandler)

	return &EchoHttpService{
		Echo: e,
	}
}
