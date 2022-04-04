package internal

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spolia/wallet-api/internal/wallet/movement"
)

type Service interface {
	CreateUser(ctx context.Context, name, lastName, alias, email string) (int64, error)
	GetBalance(ctx context.Context, alias string) (movement.AccountBalance, error)
	CreateMovement(ctx context.Context, movement movement.Movement) (int64, error)
	SearchMovement(ctx context.Context, userID int64, limit, offset uint64, movType, currencyName string) ([]movement.Row, error)
}

func API(router *gin.Engine, service Service) {
	router.POST("/users", createUser(service))
	router.GET("/users/balance/:alias", getBalance(service)) // balance
	router.POST("/movements", createMovement(service))       // send and receive
	router.GET("/movements/search", searchMovement(service))
}
