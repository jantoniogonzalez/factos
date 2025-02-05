package db

import (
	"database/sql"
	"errors"

	"github.com/jantoniogonzalez/factos/internal/models"
)

type FactosModel struct {
	database *sql.DB
}

func NewFactosModel(database *sql.DB) *FactosModel {
	return &FactosModel{database: database}
}

func (m *FactosModel) InsertOne(newFacto *models.Factos) (int, error) {
	query := `INSERT INTO factos(matchId, homeGoals, awayGoals, lastModified,
	created, userId, extraTime, penalties, homePenalties, awayPenalties, result)
	VALUES (?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?, ?, ?, ?, ?, ?);`

	args := []interface{}{
		newFacto.MatchId,
		newFacto.HomeGoals,
		newFacto.AwayGoals,
		newFacto.UserId,
		newFacto.ExtraTime,
		newFacto.Penalties,
		newFacto.HomePenalties,
		newFacto.AwayPenalties,
		newFacto.Result,
	}

	res, err := m.database.Exec(query, args...)

	if err != nil {
		return 0, err
	}

	resId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(resId), nil
}

func (m *FactosModel) GetById(id int) (*models.Factos, error) {
	query := `SELECT * FROM factos
	WHERE id=?;`

	row := m.database.QueryRow(query, id)

	f := &models.Factos{}

	err := row.Scan(&f.Id, &f.MatchId, &f.HomeGoals, &f.AwayGoals,
		&f.LastModified, &f.Created, &f.Result, &f.UserId, &f.ExtraTime, &f.Penalties)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return f, nil
}

func (m *FactosModel) GetByUser(userId int) ([]*models.Factos, error) {
	query := `SELECT * FROM factos
	WHERE userId=?;`

	rows, err := m.database.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var factos []*models.Factos

	for rows.Next() {
		f := &models.Factos{}
		if err := rows.Scan(&f.Id, &f.MatchId, &f.HomeGoals, &f.AwayGoals,
			&f.LastModified, &f.Created, &f.Result, &f.UserId, &f.ExtraTime, &f.Penalties); err != nil {
			return nil, err
		}
		factos = append(factos, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return factos, nil
}

func (m *FactosModel) Latest(quantity int) ([]*models.Factos, error) {
	query := `SELECT * FROM factos
	ORDER BY created DESC LIMIT ?;`

	rows, err := m.database.Query(query, quantity)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var factos []*models.Factos

	for rows.Next() {
		f := &models.Factos{}
		if err := rows.Scan(&f.Id, &f.MatchId, &f.HomeGoals, &f.AwayGoals,
			&f.LastModified, &f.Created, &f.Result, &f.UserId, &f.ExtraTime, &f.Penalties); err != nil {
			return nil, err
		}
		factos = append(factos, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return factos, nil

}

func (m *FactosModel) Edit(homeGoals, awayGoals, id, result int, extraTime, penalties bool) (int, error) {
	query := `UPDATE FACTOS
	SET homeGoals=?, awayGoals=?, lastModified=UTC_TIMESTAMP(), extraTime=?,
	penalties=?, result=?
	WHERE id=?;`

	res, err := m.database.Exec(query, homeGoals, awayGoals, extraTime, penalties, result, id)

	if err != nil {
		return 0, err
	}

	resId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(resId), nil
}
