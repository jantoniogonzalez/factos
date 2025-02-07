package db

import (
	"database/sql"
	"errors"

	"github.com/jantoniogonzalez/factos/internal/models"
)

type LeaguesModel struct {
	database *sql.DB
}

func NewLeaguesModel(database *sql.DB) *LeaguesModel {
	return &LeaguesModel{database: database}
}

func (m *LeaguesModel) InsertOne(newLeague models.League) (int, error) {
	query := `INSERT INTO leagues (name, apiLeagueId, country, season, logo)
	VALUES (?, ?, ?, ?, ?);`

	args := []interface{}{
		newLeague.Name,
		newLeague.ApiLeagueId,
		newLeague.Country,
		newLeague.Season,
		newLeague.Logo,
	}

	res, err := m.database.Exec(query, args...)
	if err != nil {
		// TODO: Check if the error is because apiLeagueId is not unique
		return 0, err
	}

	newId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(newId), nil
}

func (m *LeaguesModel) GetByID(leagueId int) (*models.League, error) {
	query := `SELECT * from leagues
	WHERE id=?;`

	league := &models.League{}

	err := m.database.QueryRow(query, leagueId).Scan(
		&league.ID,
		&league.Name,
		&league.ApiLeagueId,
		&league.Country,
		&league.Season,
		&league.Logo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return league, nil
}

func (m *LeaguesModel) GetByApiLeagueID(apiLeagueId int) (*models.League, error) {
	query := `SELECT * from leagues
	WHERE apiMatchId=?;`

	league := &models.League{}

	err := m.database.QueryRow(query, apiLeagueId).Scan(
		&league.ID,
		&league.Name,
		&league.ApiLeagueId,
		&league.Country,
		&league.Season,
		&league.Logo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return league, nil
}

func (m *LeaguesModel) UpdateOne(league *models.League) (int, error) {
	query := `UPDATE from leagues
	WHERE id=$
	SET name=$, country=$, season=$, logo=$;`

	args := []interface{}{
		league.ID,
		league.Name,
		league.Country,
		league.Season,
		league.Logo,
	}

	res, err := m.database.Exec(query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}

	numRowsChanged, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return int(numRowsChanged), nil
}
