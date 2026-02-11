package models

type User struct {
	Id       int64
	Email    string `binding:"required,email"`
	Password string `binding:"required,min=8"`
}
