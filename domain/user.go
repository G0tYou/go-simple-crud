package domain

import (
	"context"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserUsecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]User, string, error)
	Register(context.Context, *User) error
	Login(context.Context, string, string) (User, error)
	ChangePassword(context.Context, int64, string) error
	Delete(context.Context, int64) error
}

type UserRepository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []User, nextCursor string, err error)
	Register(ctx context.Context, u *User) error
	GetByUsername(ctx context.Context, username string) (res User, err error)
	ChangePassword(ctx context.Context, userId int64, newPassword string) error
	Delete(ctx context.Context, userId int64) error
}
