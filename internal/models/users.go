package models

import (
	"time"
)

type User struct {
	Id       int
	Username string
	GoogleId string
	Created  time.Time
}
