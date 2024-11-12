package handler

import (
	"net"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rohitxdev/abc-task/docs"
	"github.com/rohitxdev/abc-task/internal/config"
	"github.com/rohitxdev/abc-task/internal/repo"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Custom HTTP request validator
type customValidator struct {
	validator *validator.Validate
}

func (v customValidator) Validate(i any) error {
	if err := v.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	return nil
}

type Services struct {
	Config *config.Config
	Repo   *repo.Repo
}

func New(svc *Services) (*echo.Echo, error) {
	docs.SwaggerInfo.Host = net.JoinHostPort(svc.Config.Host, svc.Config.Port)

	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Validator = customValidator{
		validator: validator.New(),
	}

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		UnsafeWildcardOriginWithAllowCredentials: svc.Config.Env == "development",
	}))

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	e.POST("/classes", CreateClass(svc))
	e.POST("/bookings", CreateBooking(svc))

	return e, nil
}
