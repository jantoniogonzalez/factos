package db

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/jantoniogonzalez/factos/internal/constants"
	"github.com/jantoniogonzalez/factos/internal/models"
)

type FixturesModel struct {
	database *sql.DB
}

func NewFixturesModel(database *sql.DB) *FixturesModel {
	return &FixturesModel{database: database}
}

func (m *FixturesModel) InsertOne(newFixture models.Fixture) (int64, error) {
	_, ok := constants.MatchStatus[newFixture.MatchStatusShort]

	if !ok {
		return 0, models.ErrNoMatchingMatchStatusShort
	}

	query := `INSERT INTO fixtures(date, leagueId, homeGoals, awayGoals,
	homePenalties, awayPenalties, homeId, awayId, created, lastModified,
	matchStatusShort, apiMatchId)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?, ?);`

	res, err := m.database.Exec(
		query,
		newFixture.Date,
		newFixture.LeagueId,
		newFixture.HomeGoals,
		newFixture.AwayGoals,
		newFixture.HomePenalties,
		newFixture.AwayPenalties,
		newFixture.HomeId,
		newFixture.AwayId,
		newFixture.MatchStatusShort,
		newFixture.ApiMatchId,
	)

	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "uc_users_apiMatchId") {
				return 0, models.ErrDuplicateUsername
			}
			// if mySQLError.Number == 1000 && strings.Contains(mySQLError.Message, "") {
			// 	return 0, models.Err
			// }
		}

		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

// TODO: Add a method to insertMultiple
func (m *FixturesModel) InsertMultiple(newFixtures []*models.Fixture) (int, error) {
	return 0, nil
}

// TODO: Add a method that could be insert ignore or insert and update
func (m *FixturesModel) InsertOrUpdateMultiple(newFixtures []*models.Fixture) (int, error) {
	return 0, nil
}

func (m *FixturesModel) GetByID(fixtureId int) (*models.Fixture, error) {
	query := `SELECT * FROM fixtures
	WHERE id=?;`

	row := m.database.QueryRow(query, fixtureId)

	fixture := &models.Fixture{}

	err := row.Scan(&fixture.ID, &fixture.ApiMatchId, &fixture.Date, &fixture.LeagueId, &fixture.HomeGoals, &fixture.AwayGoals,
		&fixture.HomePenalties, &fixture.AwayPenalties, &fixture.HomeId, &fixture.AwayId, &fixture.Created, &fixture.LastModified, &fixture.MatchStatusShort)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return fixture, nil
}

func (m *FixturesModel) GetByApiMatchID(apiMatchId int) (*models.Fixture, error) {
	query := `SELECT * FROM fixtures
	WHERE apiMatchId=?;`

	row := m.database.QueryRow(query, apiMatchId)

	fixture := &models.Fixture{}

	err := row.Scan(&fixture.ID, &fixture.ApiMatchId, &fixture.Date, &fixture.LeagueId, &fixture.HomeGoals, &fixture.AwayGoals,
		&fixture.HomePenalties, &fixture.AwayPenalties, &fixture.HomeId, &fixture.AwayId, &fixture.Created, &fixture.LastModified, &fixture.MatchStatusShort)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	return fixture, nil
}

// TODO: We may want to get the teams info rather than just getting the team ids
func (m *FixturesModel) GetLatestByLeagueID(leagueId, limit int) ([]*models.Fixture, error) {
	query := `SELECT * FROM fixtures
	WHERE date<UTC_TIMESTAMP() AND leagueId=?
	ORDER BY date DESC
	LIMIT ?;`

	rows, err := m.database.Query(query, leagueId, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fixtures []*models.Fixture

	for rows.Next() {
		fixture := &models.Fixture{}
		err = rows.Scan(&fixture.ID, &fixture.ApiMatchId, &fixture.Date, &fixture.LeagueId, &fixture.HomeGoals, &fixture.AwayGoals,
			&fixture.HomePenalties, &fixture.AwayPenalties, &fixture.HomeId, &fixture.AwayId, &fixture.Created, &fixture.LastModified, &fixture.MatchStatusShort)
		if err != nil {
			return nil, err
		}
		fixtures = append(fixtures, fixture)
	}
	// Called to check if there was any error with rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return fixtures, err
}

func (m *FixturesModel) GetLatestByTeamID(teamId int, limit int) ([]*models.Fixture, error) {
	query := `SELECT * FROM fixtures
	WHERE date<UTC_TIMESTAMP() AND (homeId=? OR awayId=?)
	ORDER BY date DESC
	LIMIT ?;`

	rows, err := m.database.Query(query, teamId, teamId, limit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var fixtures []*models.Fixture

	for rows.Next() {
		fixture := &models.Fixture{}
		err = rows.Scan(&fixture.ID, &fixture.ApiMatchId, &fixture.Date, &fixture.LeagueId, &fixture.HomeGoals, &fixture.AwayGoals,
			&fixture.HomePenalties, &fixture.AwayPenalties, &fixture.HomeId, &fixture.AwayId, &fixture.Created, &fixture.LastModified, &fixture.MatchStatusShort)
		if err != nil {
			return nil, err
		}
		fixtures = append(fixtures, fixture)
	}
	// Called to check if there was any error with rows.Close()
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return fixtures, err
}

func (m *FixturesModel) UpdateByID(fixture *models.Fixture) error {
	query := `UPDATE fixtures
	SET date=$, homeGoals=$, awayGoals=$, homePenalties=$, awayPenalties=$, homeId=$, awayId=$, lastModified=UTC_TIMESTAMP(), matchStatusShort=$ 
	WHERE id=$;`

	args := []interface{}{
		fixture.Date,
		fixture.HomeGoals,
		fixture.AwayGoals,
		fixture.HomePenalties,
		fixture.AwayPenalties,
		fixture.HomeId,
		fixture.AwayId,
		fixture.MatchStatusShort,
		fixture.ID,
	}

	_, err := m.database.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *FixturesModel) UpdateByApiMatchID(fixture *models.Fixture) error {
	query := `UPDATE fixtures
	SET date=$, homeGoals=$, awayGoals=$, homePenalties=$, awayPenalties=$, homeId=$, awayId=$, lastModified=UTC_TIMESTAMP(), matchStatusShort=$ 
	WHERE apiMatchId=$;`

	args := []interface{}{
		fixture.Date,
		fixture.HomeGoals,
		fixture.AwayGoals,
		fixture.HomePenalties,
		fixture.AwayPenalties,
		fixture.HomeId,
		fixture.AwayId,
		fixture.MatchStatusShort,
		fixture.ApiMatchId,
	}

	_, err := m.database.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}
