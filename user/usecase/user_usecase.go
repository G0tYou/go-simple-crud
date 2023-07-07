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

func (u *userUsecase) Register(c context.Context, du *domain.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err = u.userRepo.Register(ctx, du)
	return
}

func (u *userUsecase) Login(c context.Context, username string, password string) (res domain.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	res, err = u.userRepo.GetByUsername(ctx, username)
	return
}

func (u *userUsecase) ChangePassword(c context.Context, userId int64, newPassword string) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err = u.userRepo.ChangePassword(ctx, userId, newPassword)
	return
}

func (u *userUsecase) Delete(c context.Context, userId int64) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	err = u.userRepo.Delete(ctx, userId)
	return
}
