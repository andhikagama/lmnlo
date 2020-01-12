package mysql

import (
	"database/sql"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/andhikagama/lmnlo/helper"
	"github.com/andhikagama/lmnlo/models/entity"
	"github.com/andhikagama/lmnlo/models/filter"
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

func (m *userRepository) Fetch(f *filter.User) ([]*entity.User, error) {
	query := sq.Select(`id, email, address`)
	query.From(`user`)

	if f.Email != `` {
		query.Where(`email = ?`, f.Email)
	}

	if f.Password != `` {
		query.Where(`password = ?`, f.Password)
	}

	if f.Address != `` {
		regx := `address REGEXP '` + f.Address + `'`
		query.Where(regx)
	}

	if f.Cursor != 0 {
		query.Where(`id  < ?`, f.Cursor)
	}

	query.OrderBy(`id DESC`).Limit(uint64(f.Num))

	query.Where(`delete_time IS NULL`)

	sql, args, _ := query.ToSql()
	res, err := m.Conn.Query(sql, args...)
	defer res.Close()

	result, err := m.unmarshal(res)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return make([]*entity.User, 0), nil
	}

	return result, err
}

func (m *userRepository) Update(usr *entity.User) (bool, error) {
	trx, err := m.Conn.Begin()

	if err != nil {
		return false, err
	}

	query := sq.Update("user").
		Set("email", usr.Email).
		Set("address", usr.Address)

	if usr.Password != `` {
		encryptedPass, _ := helper.EncryptToString(usr.Password)
		query.Set(`password`, encryptedPass)
	}

	query.Set("update_time", time.Now()).
		Where("id = ?", usr.ID)

	sql, args, _ := query.ToSql()
	stmt, err := trx.Prepare(sql)
	if err != nil {
		trx.Rollback()
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)

	if err != nil {
		trx.Rollback()
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		trx.Rollback()
		return false, err
	}

	if affected != 1 {
		trx.Rollback()
		return false, nil
	}

	usr.Password = ``
	err = trx.Commit()
	return true, nil
}

func (m *userRepository) GetByID(id int64) (*entity.User, error) {
	query := sq.Select(`id, email, address`)
	query.From(`user`)
	query.Where(`id = ?`, id)

	sql, args, _ := query.ToSql()
	res, err := m.Conn.Query(sql, args...)
	defer res.Close()

	result, err := m.unmarshal(res)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return new(entity.User), nil
	}

	return result[0], err
}

func (m *userRepository) Delete(id int64) (bool, error) {
	trx, err := m.Conn.Begin()

	if err != nil {
		return false, err
	}

	query := sq.Update("user").
		Set("delete_time", time.Now()).
		Where("id = ?", id)

	sql, args, _ := query.ToSql()

	stmt, err := trx.Prepare(sql)
	if err != nil {
		trx.Rollback()
		return false, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)

	if err != nil {
		trx.Rollback()
		return false, err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		trx.Rollback()
		return false, err
	}

	if affected != 1 {
		trx.Rollback()
		return false, nil
	}

	err = trx.Commit()
	return true, nil
}

func (m *userRepository) InsertToken(uid int64, token string) error {
	trx, err := m.Conn.Begin()
	if err != nil {
		return err
	}

	query := sq.Insert("token")
	query.Columns("user_id", "token", "create_time")
	query.Values(uid, token, time.Now())
	sql, args, _ := query.ToSql()

	stmt, err := trx.Prepare(sql)
	if err != nil {
		trx.Rollback()
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)

	if err != nil {
		trx.Rollback()
		return err
	}

	return trx.Commit()
}

func (m *userRepository) ValidateToken(token string) (bool, error) {
	query := sq.Select(`1`)
	query.From(`token`)
	query.Where(`token = ?`, token)

	sql, args, _ := query.ToSql()

	rows, err := m.Conn.Query(sql, args...)
	defer rows.Close()
	if err != nil {
		return false, err
	}

	var res int64

	for rows.Next() {

		err := rows.Scan(
			&res,
		)

		if err != nil {
			log.Error(err, res)
			return false, err
		}
	}

	if res == int64(0) {
		return false, err
	}

	return true, nil
}

func (m *userRepository) unmarshal(rows *sql.Rows) ([]*entity.User, error) {
	results := []*entity.User{}

	for rows.Next() {
		var usr entity.User

		err := rows.Scan(
			&usr.ID,
			&usr.Email,
			&usr.Address,
		)

		if err != nil {
			logrus.Error(err, usr.ID)
			return results, err
		}

		results = append(results, &usr)
	}

	return results, nil
}
