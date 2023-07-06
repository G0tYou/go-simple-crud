package mysql

import (
	"context"
	"database/sql"
	"errors"

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
	query := `SELECT id, username, password, created_at, updated_at FROM user WHERE created_at > ? ORDER BY created_at LIMIT ? `

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
