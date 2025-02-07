package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jantoniogonzalez/factos/internal/models"
)

type TeamsModel struct {
	database *sql.DB
}

func NewTeamsModel(database *sql.DB) *TeamsModel {
	return &TeamsModel{database: database}
}

func (m *TeamsModel) InsertOne(team *models.Team) (int, error) {
	query := `INSERT INTO teams(name, logo, created, lastModified, apiTeamId)
	VALUES($, $, UTC_TIMESTAMP(), UTC_TIMESTAMP(), $);`

	args := []interface{}{
		team.Name,
		team.Logo,
		team.ApiTeamId,
	}

	res, err := m.database.Exec(query, args...)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "uc_teams_apiTeamId") {
				return 0, models.ErrDuplicateApiTeamId
			}
		}
		return 0, err
	}

	teamId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(teamId), nil
}

func (m *TeamsModel) GetByID(teamId int) (*models.Team, error) {
	query := `SELECT id, name, logo FROM teams
	WHERE id=$;`

	team := &models.Team{}

	err := m.database.QueryRow(query, teamId).Scan(
		&team.ID,
		&team.Name,
		&team.Logo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return team, nil
}

func (m *TeamsModel) GetByApiTeamId(apiTeamId int) (*models.Team, error) {
	query := `SELECT id, name, logo FROM teams
	WHERE apiTeamId=$;`

	team := &models.Team{}

	err := m.database.QueryRow(query, apiTeamId).Scan(
		&team.ID,
		&team.Name,
		&team.Logo,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return team, nil
}

func (m *TeamsModel) UpdateByID(team *models.Team) (int, error) {
	query := `UPDATE teams
	WHERE id=$
	SET name=$, logo=$, lastModified=UTC_TIMESTAMP();`

	args := []interface{}{
		team.ID,
		team.Name,
		team.Logo,
	}

	res, err := m.database.Exec(query, args...)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrNoRecord
		}
		return 0, err
	}

	numTeams, err := res.RowsAffected()

	if err != nil {
		return 0, err
	}

	return int(numTeams), nil
}
