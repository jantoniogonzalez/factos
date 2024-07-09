package models

import (
	"database/sql"
	"errors"
	"time"
)

type Factos struct {
	Id           int
	MatchId      int
	LeagueId     int
	Season       int
	GoalsHome    int
	GoalsAway    int
	Result       int
	LastModified time.Time
	Created      time.Time
	UserId       int
	ExtraTime    bool
	Penalties    bool
}

type FactosModel struct {
	DB *sql.DB
}

func (m *FactosModel) Insert(matchId, leagueId, season, goalsHome, goalsAway, result, userId int, extraTime, penalties bool) (int, error) {
	query := `INSERT INTO factos(matchId, leagueId, season, goalsHome, goalsAway, lastModified,
	created, userId, extraTime, penalties, result)
	VALUES ?, ?, ?, ?, ?, UTC_TIMESTAMP(), UTC_TIMESTAMP(), ?, ?, ?, ?`

	res, err := m.DB.Exec(query, matchId, leagueId, season, goalsHome, goalsAway, userId, extraTime, penalties, result)

	if err != nil {
		return 0, err
	}

	resId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(resId), nil
}

func (m *FactosModel) GetById(id int) (*Factos, error) {
	query := `SELECT * FROM factos
	WHERE id=?`

	row := m.DB.QueryRow(query, id)

	f := &Factos{}

	err := row.Scan(&f.Id, &f.MatchId, &f.GoalsHome, &f.GoalsAway,
		&f.LastModified, &f.Created, &f.Result, &f.UserId, &f.ExtraTime, &f.Penalties)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return f, nil
}

func (m *FactosModel) GetByUser(userId int) ([]*Factos, error) {
	query := `SELECT * FROM factos
	WHERE userId=?`

	rows, err := m.DB.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var factos []*Factos

	for rows.Next() {
		f := &Factos{}
		if err := rows.Scan(&f.Id, &f.MatchId, &f.GoalsHome, &f.GoalsAway,
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

func (m *FactosModel) Latest(quantity int) ([]*Factos, error) {
	query := `SELECT * FROM factos
	ORDER BY created DESC LIMIT ?`

	rows, err := m.DB.Query(query, quantity)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var factos []*Factos

	for rows.Next() {
		f := &Factos{}
		if err := rows.Scan(&f.Id, &f.MatchId, &f.GoalsHome, &f.GoalsAway,
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

func (m *FactosModel) Edit(goalsHome, goalsAway, id, result int, extraTime, penalties bool) (int, error) {
	query := `UPDATE FACTOS
	SET goalsHome=?, goalsAway=?, lastModified=UTC_TIMESTAMP(), extraTime=?,
	penalties=?, result=?
	WHERE id=?`

	res, err := m.DB.Exec(query, goalsHome, goalsAway, extraTime, penalties, result, id)

	if err != nil {
		return 0, err
	}

	resId, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(resId), nil
}
