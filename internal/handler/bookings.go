package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rohitxdev/abc-task/internal/repo"
)

type CreateBookingRequest struct {
	MemberName string `json:"memberName" validate:"required"`
	Date       string `json:"date" validate:"required"`
	ClassID    uint64 `json:"classId" validate:"required,number"`
}

// @Summary Create a new booking
// @Description Creates a new booking for the given class and member name.
// @Tags Bookings
// @Accept json
// @Produce json
// @Param classId body uint true "ID of the class"
// @Param memberName body string true "Name of the member"
// @Param date body string true "Date of the booking in the format YYYY-MM-DD"
// @Success 200 {object} response
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /bookings [post]
func CreateBooking(svc *Services) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(CreateBookingRequest)
		if err := bindAndValidate(c, req); err != nil {
			return err
		}
		date, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "Invalid date format"})
		}
		if time.Now().After(date) {
			return c.JSON(http.StatusUnprocessableEntity, response{Message: "Date cannot be in the past"})
		}
		if err := svc.Repo.CreateBooking(c.Request().Context(), req.ClassID, req.MemberName, date.Unix()); err != nil {
			switch err {
			case repo.ClassNotFoundError:
				return c.JSON(http.StatusNotFound, response{Message: "Class not found"})
			case repo.ClassFullError:
				return c.JSON(http.StatusConflict, response{Message: "Class is full"})
			case repo.InvalidDateRangeError:
				return c.JSON(http.StatusUnprocessableEntity, response{Message: "No class is available on the given date"})
			default:
				slog.Error(err.Error())
				return echo.ErrInternalServerError
			}
		}
		return c.JSON(http.StatusCreated, response{Message: "Booking created successfully"})
	}
}
