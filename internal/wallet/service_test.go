package wallet

import (
	"context"
	"errors"
	"testing"

	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateUser_ok(t *testing.T) {
	// Given
	input := user.User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
		Password:  "1234",
	}
	// When
	var userMock userRepositoryMock
	userMock.On("Save").Return(nil).Once()
	var movementsMock movementRepositoryMock
	movementsMock.On("InitSave").Return(nil).Once()
	service := New(&userMock, &movementsMock)

	// Then
	err := service.CreateUser(context.Background(), input)
	require.NoError(t, err)
}

func TestService_CreateUser_Fail(t *testing.T) {
	// Given
	input := user.User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
		Password:  "1234",
	}
	// When
	var userMock userRepositoryMock
	userMock.On("Save").Return(errors.New("user: fail")).Once()

	service := New(&userMock, nil)

	// Then
	err := service.CreateUser(context.Background(), input)
	require.Error(t, err)
}

func TestService_CreateUser_When_InitSaveFails_DeletesUserSaved(t *testing.T) {
	// Given
	input := user.User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
		Password:  "1234",
	}
	// When
	var userMock userRepositoryMock
	userMock.On("Save").Return(nil).Once()
	userMock.On("Delete").Return(nil).Once()

	var movementsMock movementRepositoryMock
	movementsMock.On("InitSave").Return(errors.New("movement: fail")).Once()
	service := New(&userMock, &movementsMock)

	// Then
	err := service.CreateUser(context.Background(), input)
	require.Error(t, err)
}

func TestService_GetUser_When_GetAccountExtractFail_Then_ReturnsError(t *testing.T) {
	// When
	var userMock userRepositoryMock
	userMock.On("Get").Return(user.User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
	}, nil).Once()
	userMock.On("Delete").Return(nil).Once()

	var movementsMock movementRepositoryMock
	movementsMock.On("GetAccountExtract").Return(movement.AccountBalance{}, errors.New("mov fail")).Once()
	service := New(&userMock, &movementsMock)

	// Then
	userResult, err := service.GetBalance(context.Background(), "user")
	require.Error(t, err)
	require.Empty(t, userResult)
}

func TestService_Send_ok(t *testing.T) {
	// Given
	input := movement.Movement{
		Type:             "deposit",
		Amount:           100,
		CurrencyName:     "ARS",
		Alias:            "user",
		InteractionAlias: "user",
	}
	// When
	var userMock userRepositoryMock
	var movementsMock movementRepositoryMock
	userMock.On("Exist").Return(true, nil).Once()
	movementsMock.On("Save").Return(nil).Once()
	movementsMock.On("GetFunds").Return(float64(100), nil).Once()
	service := New(&userMock, &movementsMock)

	// Then
	err := service.Send(context.Background(), input)
	require.NoError(t, err)
}

func TestService_Send_Fail(t *testing.T) {
	// Given
	input := movement.Movement{
		Type:         "deposit",
		Amount:       100,
		CurrencyName: "ARS",
		Alias:        "user",
	}
	// When
	var userMock userRepositoryMock
	var movementsMock movementRepositoryMock
	userMock.On("Exist").Return(true, nil).Once()
	movementsMock.On("GetFunds").Return(float64(100), nil).Once()
	movementsMock.On("Save").Return(errors.New("movement:fail")).Once()
	service := New(&userMock, &movementsMock)

	// Then
	err := service.Send(context.Background(), input)
	require.Error(t, err)
}

type userRepositoryMock struct {
	mock.Mock
}

type movementRepositoryMock struct {
	mock.Mock
}

func (u *userRepositoryMock) Save(ctx context.Context, us user.User) error {
	args := u.Called()
	return args.Error(0)
}

func (u *userRepositoryMock) Exist(ctx context.Context, alias string) (bool, error) {
	args := u.Called()
	return args.Bool(0), args.Error(1)
}

func (u *userRepositoryMock) Delete(ctx context.Context, alias string) error {
	args := u.Called()
	return args.Error(0)
}

func (u *userRepositoryMock) IsValidCredential(ctx context.Context, alias, password string) (bool, error) {
	args := u.Called()
	return args.Bool(0), args.Error(1)
}

func (m *movementRepositoryMock) Save(ctx context.Context, movement movement.Movement) error {
	args := m.Called()
	return args.Error(0)
}

func (m *movementRepositoryMock) InitSave(ctx context.Context, movement movement.Movement) error {
	args := m.Called()
	return args.Error(0)
}

func (m *movementRepositoryMock) GetHistory(ctx context.Context, alias string) (movement.AccountHistory, error) {
	args := m.Called()
	return args.Get(0).(movement.AccountHistory), args.Error(1)
}

func (m *movementRepositoryMock) GetFunds(ctx context.Context, currencyName, alias string) (float64, error) {
	args := m.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (m *movementRepositoryMock) GetAccountExtract(ctx context.Context, alias string) (movement.AccountBalance, error) {
	args := m.Called()
	return args.Get(0).(movement.AccountBalance), args.Error(1)
}
