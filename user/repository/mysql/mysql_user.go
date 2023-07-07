package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"simple_crud/domain"
	"simple_crud/user/repository"

	"github.com/sirupsen/logrus"
)

type mysqlUserRepository struct {
	Conn *sql.DB
}

func NewMysqlUserRepository(conn *sql.DB) domain.UserRepository {
	return &mysqlUserRepository{conn}
}

func (m *mysqlUserRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.User, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.User, 0)
	for rows.Next() {
		t := domain.User{}
		err = rows.Scan(
			&t.ID,
			&t.Username,
			&t.Password,
			&t.Active,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlUserRepository) Fetch(ctx context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	query := `SELECT id, username, password, active,created_at, updated_at FROM user WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodeCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", errors.New("Param is not valid")
	}

	res, err = m.fetch(ctx, query, decodeCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}

func (m *mysqlUserRepository) Register(ctx context.Context, du *domain.User) (err error) {
	query := "INSERT user SET username = ?, password = ?, created_at = ?, updated_at = ?"
	pass, err := repository.EncodePassword(du.Password)
	if err != nil {
		return
	}
	du.Password = pass

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, du.Username, du.Password, du.CreatedAt, du.UpdatedAt)
	if err != nil {
		return
	}

	id, err := res.LastInsertId()
	du.ID = id
	return
}

func (m *mysqlUserRepository) GetByUsername(ctx context.Context, username string) (res domain.User, err error) {
	query := "SELECT id, username, password, active,created_at, updated_at From user WHERE username = ?"

	list, err := m.fetch(ctx, query, username)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrUserNotFound
	}

	return
}

func (m *mysqlUserRepository) ChangePassword(ctx context.Context, userId int64, newPassword string) (err error) {
	query := "UPDATE user SET  password = ? WHERE id = ?"
	newPassword, err = repository.EncodePassword(newPassword)
	if err != nil {
		return
	}

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, newPassword, userId)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}

func (m *mysqlUserRepository) Delete(ctx context.Context, userId int64) (err error) {
	query := "DELETE FROM user WHERE id = ?"
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, userId)
	if err != nil {
		return
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
