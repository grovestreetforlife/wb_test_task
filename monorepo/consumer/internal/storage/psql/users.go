package psql

import (
	"context"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"wb_test_task/consumer/internal/common"
	"wb_test_task/consumer/internal/domain"
	"wb_test_task/libs/model"
)

type userStorage struct {
	pool pool
}

func newUserStorage(pool pool) *userStorage {
	return &userStorage{pool: pool}
}

// Create создание пользователя
func (u *userStorage) Create(ctx context.Context, id string) (*model.User, error) {
	query := `
		INSERT INTO users(id) VALUES ($1)
	`

	_, err := u.pool.Exec(ctx, query, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case domain.CodeErrDuplicateKey:
				return &model.User{}, common.WrapError{Err: domain.ErrUserAlreadyExists, Msg: domain.ErrUserAlreadyExists.Error()}
			case domain.CodeErrConstraintLenValue:
				return &model.User{}, common.WrapError{Err: domain.ErrInvalidUserID, Msg: "length of id bigger them 500"}
			}
		}
		return &model.User{}, common.WrapError{Err: err, Msg: "fail to create user"}
	}

	return &model.User{ID: id}, nil
}

func (u *userStorage) GetByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT u.id FROM users u WHERE id=$1
	`

	var user model.User
	if err := u.pool.QueryRow(ctx, query, id).Scan(&user.ID); err != nil {
		return &model.User{}, err
	}

	return &user, nil
}
