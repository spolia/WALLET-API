package user

import (
	"context"
	"errors"
)

var (
	ErrorInvalidCredential   = errors.New("invalid credentials")
	ErrorDestinyUserNotFound = errors.New("not found")
	ErrorCreatingUser        = errors.New("creating user")
	ErrorAlreadyExist        = errors.New("already exist")
)

type Repository interface {
	Save(ctx context.Context, u User) error
	Exist(ctx context.Context, alias string) (bool, error)
	Delete(ctx context.Context, alias string) error
	IsValidCredential(ctx context.Context, alias, password string) (bool, error)
}

type User struct {
	Alias           string             `json:"alias" validate:"required"`
	FirstName       string             `json:"firstname" validate:"required"`
	LastName        string             `json:"lastname" validate:"required"`
	Email           string             `json:"email" validate:"required"`
	WalletStatement map[string]float64 `json:"walletstatement"`
	Password        string             `json:"password" validate:"required"`
}
