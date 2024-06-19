package models

import "time"

type User struct {
	Id       int
	Username string
	Created  time.Time
}
