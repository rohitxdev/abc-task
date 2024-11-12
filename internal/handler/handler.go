package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/abc-task/internal/config"
	"github.com/rohitxdev/abc-task/internal/repo"
)

// Generic response
type response struct {
	Message string `json:"message"`
}

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
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true

	e.Validator = customValidator{
		validator: validator.New(),
	}

	e.POST("/classes", CreateClass(svc))

	e.POST("/bookings", CreateBooking(svc))

	return e, nil
}

// bindAndValidate binds path params, query params and the request body into provided type `i` and validates provided `i`. `i` must be a pointer. The default binder binds body based on Content-Type header. Validator must be registered using `Echo#Validator`.
func bindAndValidate(c echo.Context, i any) error {
	var err error
	if err = c.Bind(i); err != nil {
		return err
	}
	binder := echo.DefaultBinder{}
	if err = binder.BindHeaders(c, i); err != nil {
		return err
	}
	if err = c.Validate(i); err != nil {
		return echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}
	return err
}
