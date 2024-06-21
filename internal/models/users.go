package models

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Id       int
	Username string
	GoogleId string
	Created  time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(username, googleId string) (int, error) {
	query := `INSERT INTO Users values(username, googleId, created)
	values(?, ?, UTC_TIMESTAMP());`

	res, err := m.DB.Exec(query, username, googleId)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Get(googleId string) (*User, error) {
	query := `SELECT id, username, created FROM users
	WHERE googleId=?;`

	row := m.DB.QueryRow(query, googleId)

	user := &User{}

	err := row.Scan(&user.Id, &user.Username, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return user, nil
}
