package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id           int
	Username     string
	Created      time.Time
	RefreshToken string
	AccessToken  string
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username string) error {
	// query := `INSERT INTO Users values(username, created, refreshToken, accessToken)
	// values()`

	return nil
}
