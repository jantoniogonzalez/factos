package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jantoniogonzalez/factos/internal/models"
)

type UserModel struct {
	database *sql.DB
}

func NewUserModel(database *sql.DB) *UserModel {
	return &UserModel{database: database}
}

func (m *UserModel) Insert(username, googleId string) (int, error) {
	query := `INSERT INTO users (username, googleId, created)
	VALUES (?, ?, UTC_TIMESTAMP());`

	result, err := m.database.Exec(query, username, googleId)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "uc_users_username") {
				return 0, models.ErrDuplicateUsername
			}
		}
		return 0, err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *UserModel) Get(googleId string) (*models.User, error) {
	query := `SELECT id, username, created FROM users
	WHERE googleId=?;`

	row := m.database.QueryRow(query, googleId)

	user := &models.User{}

	err := row.Scan(&user.Id, &user.Username, &user.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return user, nil
}
