package usecase

import (
	"context"
	"simple_crud/domain"
	"time"
)

type userUsecase struct {
	userRepo       domain.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(u domain.UserRepository, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepo:       u,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Fetch(c context.Context, cursor string, num int64) (res []domain.User, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, nextCursor, err = u.userRepo.Fetch(ctx, cursor, num)

	if err != nil {
		return nil, "", err
	}

	return
}
