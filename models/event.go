package models

import (
	"time"
)

type Event struct {
	Id          int64
	Name        string    `binding:"required,min=3,max=100"`
	Description string    `binding:"required,min=5,max=500"`
	Location    string    `binding:"required,min=3,max=100"`
	DateTime    time.Time `binding:"required,futuredate"`
	UserId      int64
}
