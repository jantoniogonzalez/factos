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

func (m *FixturesModel) GetFutureByLeague(leagueId int, limit int) ([]*models.Fixture, error) {
	query := `SELECT id, date, leagueId, homeGoals, awayGoals, homePenalties, awayPenalties, 
              homeId, awayId, created, lastModified, matchStatusShort, apiMatchId
    FROM fixtures
    WHERE leagueId = ? AND date > UTC_TIMESTAMP()
    ORDER BY date ASC
    LIMIT ?`

	rows, err := m.database.Query(query, leagueId, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fixtures []*models.Fixture
	for rows.Next() {
		f := &models.Fixture{}
		err = rows.Scan(&f.ID, &f.Date, &f.LeagueId, &f.HomeGoals, &f.AwayGoals, &f.HomePenalties,
			&f.AwayPenalties, &f.HomeId, &f.AwayId, &f.Created, &f.LastModified, &f.MatchStatusShort, &f.ApiMatchId)
		if err != nil {
			return nil, err
		}
		fixtures = append(fixtures, f)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return fixtures, nil
}
