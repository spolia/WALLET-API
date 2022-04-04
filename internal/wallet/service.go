package wallet

import (
	"context"
	"strings"

	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

type Service struct {
	userRepo     user.Repository
	movementRepo movement.Repository
}

// New creates a Service implementation.
func New(userRepo user.Repository, movRepo movement.Repository) *Service {
	return &Service{userRepo: userRepo, movementRepo: movRepo}
}

// CreateUser saves a new user
func (s *Service) CreateUser(ctx context.Context, name, lastName, alias, email string) (int64, error) {
	userID, err := s.userRepo.Save(ctx, name, lastName, alias, email)
	if err != nil {
		println(err.Error())
		return 0, err
	}

	// every time that a new user is saved is necessary init movements
	err = s.movementRepo.InitSave(ctx, movement.Movement{
		Type:             "init",
		Alias:            alias,
		InteractionAlias: alias,
	})
	if err != nil {
		// delete the created user
		println("Delete", err.Error())
		return 0, s.userRepo.Delete(ctx, userID)
	}

	return userID, nil
}

// GetBalance returns an balance user
func (s *Service) GetBalance(ctx context.Context, alias string) (movement.AccountBalance, error) {
	accountExtract, err := s.movementRepo.GetAccountExtract(ctx, alias)
	if err != nil {
		return movement.AccountBalance{}, err
	}

	return accountExtract, nil
}

// CreateMovement saves a movement
func (s *Service) CreateMovement(ctx context.Context, movement movement.Movement) (int64, error) {
	movement.CurrencyName = strings.ToUpper(movement.CurrencyName)
	movementID, err := s.movementRepo.Save(ctx, movement)
	if err != nil {
		return 0, err
	}

	return movementID, nil
}
