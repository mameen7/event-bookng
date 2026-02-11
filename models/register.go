package models

type RegisterEvent struct {
	Id      int64
	UserId  int64 `binding:"required"`
	EventId int64 `binding:"required"`
}
