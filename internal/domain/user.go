package domain

import "time"

type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string
	Birthday time.Time
	Gender   int8
	Phone    string
	Profile  string
}
