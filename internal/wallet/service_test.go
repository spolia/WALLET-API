package wallet

import (
	"context"
	"errors"
	"testing"
	"time"

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
	}
	// When
	var userMock userRepositoryMock
	userMock.On("Save").Return(int64(1), nil).Once()
	var movementsMock movementRepositoryMock
	movementsMock.On("InitSave").Return(nil).Once()
	service := New(&userMock, &movementsMock)

	// Then
	userID, err := service.CreateUser(context.Background(), input.FirstName, input.LastName, input.Alias, input.Email)
	require.NoError(t, err)
	require.Equal(t, int64(1), userID)
}

func TestService_CreateUser_Fail(t *testing.T) {
	// Given
	input := user.User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
	}
	// When
	var userMock userRepositoryMock
	userMock.On("Save").Return(int64(0), errors.New("user: fail")).Once()

	service := New(&userMock, nil)

	// Then
	userID, err := service.CreateUser(context.Background(), input.FirstName, input.LastName, input.Alias, input.Email)
	require.Error(t, err)
	require.Equal(t, int64(0), userID)
}

func TestService_CreateUser_When_InitSaveFails_DeletesUserSaved(t *testing.T) {
	// Given
	input := user.User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
	}
	// When
	var userMock userRepositoryMock
	userMock.On("Save").Return(int64(1), nil).Once()
	userMock.On("Delete").Return(nil).Once()

	var movementsMock movementRepositoryMock
	movementsMock.On("InitSave").Return(errors.New("movement: fail")).Once()
	service := New(&userMock, &movementsMock)

	// Then
	userID, err := service.CreateUser(context.Background(), input.FirstName, input.LastName, input.Alias, input.Email)
	require.NoError(t, err)
	require.Equal(t, int64(0), userID)
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
	userResult, err := service.GetUser(context.Background(), 1)
	require.Error(t, err)
	require.Empty(t, userResult)
}

func TestService_CreateMovement_ok(t *testing.T) {
	// Given
	input := movement.Movement{
		Type:         "deposit",
		Amount:       100,
		CurrencyName: "ARS",
		UserID:       1,
	}
	// When
	var userMock userRepositoryMock
	var movementsMock movementRepositoryMock
	movementsMock.On("Save").Return(int64(1), nil).Once()
	service := New(&userMock, &movementsMock)

	// Then
	id, err := service.CreateMovement(context.Background(), input)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)
}

func TestService_CreateMovement_Fail(t *testing.T) {
	// Given
	input := movement.Movement{
		Type:         "deposit",
		Amount:       100,
		CurrencyName: "ARS",
		UserID:       1,
	}
	// When
	var userMock userRepositoryMock
	var movementsMock movementRepositoryMock
	movementsMock.On("Save").Return(int64(0), errors.New("movement:fail")).Once()
	service := New(&userMock, &movementsMock)

	// Then
	id, err := service.CreateMovement(context.Background(), input)
	require.Error(t, err)
	require.Equal(t, int64(0), id)
}

func TestService_SearchMovement_Ok(t *testing.T) {
	// When
	var userMock userRepositoryMock
	var movementsMock movementRepositoryMock
	movementsMock.On("Search").Return([]movement.Row{
		{
			CurrencyName: "USDT",
			Type:         "deposut",
			DateCreated:  time.Now(),
			Amount:       100.00,
			TotalAmount:  200.00,
		},
		{
			CurrencyName: "USDT",
			Type:         "deposut",
			DateCreated:  time.Now(),
			Amount:       100.00,
			TotalAmount:  300.00,
		},
	}, nil).Once()
	service := New(&userMock, &movementsMock)

	// Then
	movements, err := service.SearchMovement(context.Background(), 1, uint64(10), uint64(0), "deposit",
		"usdt")
	require.NoError(t, err)
	require.Equal(t, 2, len(movements))
	require.Equal(t, 200.00, movements[0].TotalAmount)
}

func TestService_SearchMovement_Fail(t *testing.T) {
	// When
	var userMock userRepositoryMock
	var movementsMock movementRepositoryMock
	movementsMock.On("Search").Return([]movement.Row{}, errors.New("fail")).Once()
	service := New(&userMock, &movementsMock)

	// Then
	movements, err := service.SearchMovement(context.Background(), 1, uint64(10), uint64(0), "deposit",
		"usdt")
	require.Error(t, err)
	require.Equal(t, 0, len(movements))
}

type userRepositoryMock struct {
	mock.Mock
}

type movementRepositoryMock struct {
	mock.Mock
}

func (u *userRepositoryMock) Save(ctx context.Context, firstName, lastName, alias, email string) (int64, error) {
	args := u.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (u *userRepositoryMock) Get(ctx context.Context, id int64) (user.User, error) {
	args := u.Called()
	return args.Get(0).(user.User), args.Error(1)
}

func (u *userRepositoryMock) Delete(ctx context.Context, id int64) error {
	args := u.Called()
	return args.Error(0)
}

func (m *movementRepositoryMock) Save(ctx context.Context, movement movement.Movement) (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (m *movementRepositoryMock) InitSave(ctx context.Context, movement movement.Movement) error {
	args := m.Called()
	return args.Error(0)
}

func (m *movementRepositoryMock) Search(ctx context.Context, userID int64, limit, offset uint64, movType,
	currencyName string) ([]movement.Row, error) {
	args := m.Called()
	return args.Get(0).([]movement.Row), args.Error(1)
}

func (m *movementRepositoryMock) GetAccountExtract(ctx context.Context, id int64) (movement.AccountBalance, error) {
	args := m.Called()
	return args.Get(0).(movement.AccountBalance), args.Error(1)
}
