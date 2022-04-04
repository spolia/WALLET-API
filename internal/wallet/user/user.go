package user

import (
	"context"
	"errors"
)

var ErrorUserNotFound = errors.New("user: not found")
var ErrorAlreadyExist = errors.New("user: already exist")

type Repository interface {
	Save(ctx context.Context, firstName, lastName, alias, email string) (int64, error)
	Get(ctx context.Context, id int64) (User, error)
	Delete(ctx context.Context, id int64) error
}

type User struct {
	ID              int64              `json:"id"`
	FirstName       string             `json:"firstname" binding:"required"`
	LastName        string             `json:"lastname" binding:"required"`
	Alias           string             `json:"alias" binding:"required"`
	Email           string             `json:"email" binding:"required"`
	WalletStatement map[string]float64 `json:"walletstatement"`
	Password        string             `json:"password"`
}
