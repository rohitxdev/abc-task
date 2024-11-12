package repo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/rohitxdev/abc-task/internal/database"
	"github.com/rohitxdev/abc-task/internal/repo"
	"github.com/stretchr/testify/assert"
)

func TestRepo(t *testing.T) {
	db, err := database.NewSQLite("test.db")
	assert.Nil(t, err)
	defer func() {
		db.Close()
		assert.Nil(t, os.RemoveAll(database.DirName))
	}()

	r, err := repo.New(db)
	assert.Nil(t, err)

	t.Run("CreateClass", func(t *testing.T) {
		type args struct {
			name      string
			startDate int64
			endDate   int64
			capacity  uint
		}
		tests := []struct {
			name string
			args args
			want error
		}{
			{
				name: "Valid args",
				args: args{
					name:      "Yoga-1",
					startDate: time.Now().Add(time.Hour * 24).Unix(),
					endDate:   time.Now().Add(time.Hour * 24 * 2).Unix(),
					capacity:  3,
				},
				want: nil,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := r.CreateClass(context.TODO(), tt.args.name, tt.args.startDate, tt.args.endDate, tt.args.capacity)
				assert.Equal(t, tt.want, err)
			})
		}
	})

	t.Run("CreateBooking", func(t *testing.T) {
		type args struct {
			memberName string
			date       int64
			classID    uint64
		}
		tests := []struct {
			name string
			args args
			want error
		}{
			{
				name: "Valid args",
				args: args{
					memberName: "Rohit",
					date:       time.Now().Add(time.Hour * 24).Unix(),
					classID:    1,
				},
				want: nil,
			},
			{
				name: "Invalid classID",
				args: args{
					memberName: "Rohit",
					date:       time.Now().Add(time.Hour * 24).Unix(),
					classID:    0,
				},
				want: repo.ClassNotFoundError,
			},
			{
				name: "Invalid date range",
				args: args{
					memberName: "Rohit",
					date:       time.Now().Add(time.Hour * 24 * 100).Unix(),
					classID:    1,
				},
				want: repo.InvalidDateRangeError,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := r.CreateBooking(context.TODO(), tt.args.classID, tt.args.memberName, tt.args.date)
				assert.Equal(t, tt.want, err)
			})
		}

	})
}
