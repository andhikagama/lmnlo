package mysql

import (
	"database/sql"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/user"
	sq "github.com/elgris/sqrl"
	"github.com/sirupsen/logrus"
)

type userRepository struct {
	Conn *sql.DB
}

// NewUserRepository ...
func NewUserRepository(Conn *sql.DB) user.Repository {
	return &userRepository{Conn}
}

func (m *userRepository) Store(usr *entity.User) error {
	trx, err := m.Conn.Begin()
	if err != nil {
		return err
	}

	query := sq.Insert(`user`)
	query.Columns(`email`, `password`, `address`, `create_time`)
	query.Values(usr.Email, usr.Password, usr.Address, time.Now())

	sql, args, _ := query.ToSql()

	stmt, err := trx.Prepare(sql)
	if err != nil {
		trx.Rollback()
		return err
	}
	defer stmt.Close()

	r, err := stmt.Exec(args...)

	if err != nil {
		trx.Rollback()
		return err
	}

	var id int64

	id, err = r.LastInsertId()
	if err != nil {
		log.Error(err, id)
		trx.Rollback()
		return err
	}

	usr.ID = id

	return trx.Commit()
}
func (m *userRepository) unmarshal(rows *sql.Rows) ([]*entity.User, error) {
	results := []*entity.User{}

	for rows.Next() {
		var usr entity.User

		err := rows.Scan(
			&usr.ID,
		)

		if err != nil {
			logrus.Error(err, usr.ID)
			return results, err
		}

		results = append(results, &usr)
	}

	return results, nil
}
