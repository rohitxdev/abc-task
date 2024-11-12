package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/rohitxdev/abc-task/internal/config"
	"github.com/rohitxdev/abc-task/internal/database"
	"github.com/rohitxdev/abc-task/internal/handler"
	"github.com/rohitxdev/abc-task/internal/repo"
	"github.com/stretchr/testify/assert"
)

type httpRequestOpts struct {
	query   map[string]string
	body    any
	headers map[string]string
	method  string
	path    string
}

func createHttpRequest(opts *httpRequestOpts) (*http.Request, error) {
	url, err := url.Parse(opts.path)
	if err != nil {
		return nil, err
	}
	q := url.Query()
	for key, value := range opts.query {
		q.Set(key, value)
	}
	url.RawQuery = q.Encode()
	j, err := json.Marshal(opts.body)
	if err != nil {
		return nil, err
	}
	req := httptest.NewRequest(opts.method, url.String(), bytes.NewReader(j))
	for key, value := range opts.headers {
		req.Header.Set(key, value)
	}
	return req, err
}

func TestAPI(t *testing.T) {
	cfg, err := config.Load()
	assert.Nil(t, err)

	db, err := database.NewSQLite("test.db")
	assert.Nil(t, err)
	defer func() {
		db.Close()
		assert.Nil(t, os.RemoveAll(database.DirName))
	}()

	r, err := repo.New(db)
	assert.Nil(t, err)

	svc := &handler.Services{
		Config: cfg,
		Repo:   r,
	}

	h, err := handler.New(svc)
	assert.Nil(t, err)

	t.Run("POST /classes", func(t *testing.T) {
		type args struct {
			body handler.CreateClassRequest
		}
		tests := []struct {
			name string
			args args
			want int
		}{
			{name: "Valid request", args: args{
				body: handler.CreateClassRequest{
					Name:      "Yoga-1",
					StartDate: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					EndDate:   time.Now().Add(time.Hour * 24 * 2).Format("2006-01-02"),
					Capacity:  3,
				}},
				want: http.StatusCreated,
			},
			{name: "Start date is in the past", args: args{
				body: handler.CreateClassRequest{
					Name:      "Yoga-2",
					StartDate: time.Now().Add(time.Hour * -24).Format("2006-01-02"),
					EndDate:   time.Now().Add(time.Hour * 24 * 2).Format("2006-01-02"),
					Capacity:  3,
				}},
				want: http.StatusUnprocessableEntity,
			},
			{name: "End date is in the past", args: args{
				body: handler.CreateClassRequest{
					Name:      "Yoga-3",
					StartDate: time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					EndDate:   time.Now().Add(time.Hour * -24).Format("2006-01-02"),
					Capacity:  3,
				}},
				want: http.StatusUnprocessableEntity,
			},
			{name: "End date is before start date", args: args{
				body: handler.CreateClassRequest{
					Name:      "Yoga-4",
					StartDate: time.Now().Add(time.Hour * 24 * 2).Format("2006-01-02"),
					EndDate:   time.Now().Add(time.Hour * 24).Format("2006-01-02"),
					Capacity:  3,
				}},
				want: http.StatusUnprocessableEntity,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := createHttpRequest(&httpRequestOpts{
					method: http.MethodPost,
					path:   "/classes",
					body:   tt.args.body,
					headers: map[string]string{
						"Content-Type": "application/json",
					},
				})
				assert.Nil(t, err)
				res := httptest.NewRecorder()
				c := h.NewContext(req, res)
				err = handler.CreateClass(svc)(c)
				assert.Nil(t, err)
				assert.Equal(t, tt.want, res.Code)
			})
		}
	})

	t.Run("POST /bookings", func(t *testing.T) {
		type args struct {
			body handler.CreateBookingRequest
		}
		tests := []struct {
			name string
			args args
			want int
		}{
			{
				name: "Class not available on the given date",
				args: args{
					body: handler.CreateBookingRequest{
						ClassID:    1,
						MemberName: "Rohit",
						Date:       time.Now().Add(time.Hour * 24 * 100).Format("2006-01-02"),
					},
				},
				want: http.StatusUnprocessableEntity,
			},
			{name: "Date is in the past", args: args{
				body: handler.CreateBookingRequest{
					ClassID:    1,
					MemberName: "Rohit",
					Date:       time.Now().Add(time.Hour * -24).Format("2006-01-02"),
				}},
				want: http.StatusUnprocessableEntity,
			},
			{name: "Class not found", args: args{
				body: handler.CreateBookingRequest{
					ClassID:    2,
					MemberName: "Rohit",
					Date:       time.Now().Add(time.Hour * 24).Format("2006-01-02"),
				}},
				want: http.StatusNotFound,
			},
			{name: "Valid request", args: args{
				body: handler.CreateBookingRequest{
					ClassID:    1,
					MemberName: "Rohit",
					Date:       time.Now().Add(time.Hour * 24).Format("2006-01-02"),
				}},
				want: http.StatusCreated,
			},
			{name: "Valid request", args: args{
				body: handler.CreateBookingRequest{
					ClassID:    1,
					MemberName: "Rohit",
					Date:       time.Now().Add(time.Hour * 24).Format("2006-01-02"),
				}},
				want: http.StatusCreated,
			},
			{name: "Valid request", args: args{
				body: handler.CreateBookingRequest{
					ClassID:    1,
					MemberName: "Rohit",
					Date:       time.Now().Add(time.Hour * 24).Format("2006-01-02"),
				}},
				want: http.StatusCreated,
			},
			{name: "Class is full", args: args{
				body: handler.CreateBookingRequest{
					ClassID:    1,
					MemberName: "Rohit",
					Date:       time.Now().Add(time.Hour * 24).Format("2006-01-02"),
				}},
				want: http.StatusConflict,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := createHttpRequest(&httpRequestOpts{
					method: http.MethodPost,
					path:   "/bookings",
					body:   tt.args.body,
					headers: map[string]string{
						"Content-Type": "application/json",
					},
				})
				assert.Nil(t, err)
				res := httptest.NewRecorder()
				c := h.NewContext(req, res)
				err = handler.CreateBooking(svc)(c)
				assert.Nil(t, err)
				assert.Equal(t, tt.want, res.Code)
			})
		}
	})
}
