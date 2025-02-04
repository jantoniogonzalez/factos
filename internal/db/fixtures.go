package db

import (
	"database/sql"
)

type FixturesModel struct {
	DB *sql.DB
}

func (FixturesModel) Insert() (int, error) {
	return -1, nil
}

func (FixturesModel) GetFixtureByID(fixtureId int) {
	//query

	//m.DB.Exec

}
