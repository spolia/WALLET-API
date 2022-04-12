package wallet

import (
	"context"

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
func (s *Service) CreateUser(ctx context.Context, u user.User) error {
	err := s.userRepo.Save(ctx, u)
	if err != nil {
		println(err.Error())
		return err
	}
	// every time that a new user is saved is necessary init account movements
	err = s.movementRepo.InitSave(ctx, movement.Movement{
		Type:             "init",
		Alias:            u.Alias,
		InteractionAlias: u.Alias,
	})

	if err != nil {
		// delete the created user
		s.userRepo.Delete(ctx, u.Alias)
		return user.ErrorCreatingUser
	}

	return nil
}

// GetBalance returns a balance user
func (s *Service) GetBalance(ctx context.Context, alias string) (movement.AccountBalance, error) {
	accountExtract, err := s.movementRepo.GetAccountExtract(ctx, alias)
	if err != nil {
		return movement.AccountBalance{}, err
	}

	return accountExtract, nil
}

// Send the money to other user account if the user have is funds
// otherwise returns error
func (s *Service) Send(ctx context.Context, m movement.Movement) error {
	ok, err := s.userRepo.Exist(ctx, m.InteractionAlias)
	if err != nil {
		return err
	}
	if !ok {
		return user.ErrorDestinyUserNotFound
	}

	// check the funds
	funds, err := s.movementRepo.GetFunds(ctx, m.CurrencyName, m.Alias)
	if funds == 0 || funds-m.TotalAmount < 0 {
		return movement.ErrorInsufficientFunds
	}

	err = s.movementRepo.Save(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

// AutoDeposit deposit money into the user account
func (s *Service) AutoDeposit(ctx context.Context, m movement.Movement) error {
	err := s.movementRepo.Save(ctx, m)
	if err != nil {
		return err
	}

	return nil
}

// GetHistory returns the movements history for all user accounts
func (s *Service) GetHistory(ctx context.Context, alias string) (movement.AccountHistory, error) {
	history, err := s.movementRepo.GetHistory(ctx, alias)
	if err != nil {
		return movement.AccountHistory{}, err
	}

	return history, nil
}

// ValidateCredential given a alias and passwords returns true if exist the user
// otherwise returns error
func (s *Service) ValidateCredential(ctx context.Context, alias, password string) (bool, error) {
	return s.userRepo.IsValidCredential(ctx, alias, password)
}
