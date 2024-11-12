package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type CreateClassRequest struct {
	Name      string `json:"name" validate:"required"`
	StartDate string `json:"startDate" validate:"required"`
	EndDate   string `json:"endDate" validate:"required"`
	Capacity  uint   `json:"capacity" validate:"required"`
}

// @Summary Create a new class
// @Description Creates a new class with the given name, description, start date, end date, and capacity.
// @Tags Classes
// @Accept json
// @Produce json
// @Param body body handler.CreateClassRequest true "Request body"
// @Success 201 {object} response
// @Failure 400 {object} response
// @Failure 422 {object} response
// @Failure 500 {object} response
// @Router /classes [post]
func CreateClass(svc *Services) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(CreateClassRequest)
		if err := bindAndValidate(c, req); err != nil {
			return err
		}

		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "Invalid date format for start date"})
		}
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "Invalid date format for end date"})
		}

		t := time.Now()
		if t.After(startDate) {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "Start date cannot be in the past"})
		}
		if t.After(endDate) {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "End date cannot be in the past"})
		}
		if startDate.After(endDate) {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "End date cannot be before start date"})
		}
		if err := svc.Repo.CreateClass(req.Name, startDate.Unix(), endDate.Unix(), req.Capacity); err != nil {
			switch err {
			default:
				// Usually I add a lot more details to the log for internal server errors, but for this task, I'm just logging the error and returning a generic error message
				slog.Error(err.Error())
				return echo.ErrInternalServerError
			}
		}

		return c.JSON(http.StatusCreated, response{Message: "Class created successfully"})
	}
}
