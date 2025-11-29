package models

import (
	"time"
)

type User struct {
	Id           int       `json:"id"`
	Username     string    `json:"username"`
	GoogleId     string    `json:"googleId"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}
