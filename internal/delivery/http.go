package delivery

import (
	"context"
	"net/http"

	_ "github.com/arvaliullin/wapa/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
)

type EchoHttpService struct {
	Echo *echo.Echo
}

func (e *EchoHttpService) Start(address string) error {
	return e.Echo.Start(address)
}

// Shutdown останавливает HTTP сервер с учётом контекста.
func (e *EchoHttpService) Shutdown(ctx context.Context) error {
	return e.Echo.Shutdown(ctx)
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
