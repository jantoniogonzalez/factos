package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
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

func (m *UserModel) Insert(username, googleId string) error {
	query := `INSERT INTO Users values(username, googleId, created)
	values(?, ?, UTC_TIMESTAMP());`

	_, err := m.DB.Exec(query, username, googleId)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "uc_users_username") {
				return ErrDuplicateUsername
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) Get(googleId string) (*User, error) {
	query := `SELECT id, username, created FROM users
	WHERE googleId=?;`

	row := m.DB.QueryRow(query, googleId)

	user := &User{}

	err := row.Scan(&user.Id, &user.Username, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}
