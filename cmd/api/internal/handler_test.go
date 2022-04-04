package internal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

func Test_Handler_API_createUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tt := []struct {
		TestName, Filename string
		ExpectedStatus     int
		Error              error
	}{
		{"Ok", "create_user_ok", http.StatusCreated, nil},
		{"WrongFormat", "create_user_wrong_format", http.StatusBadRequest, nil},
		{"ErrorAlreadyExist", "create_user_ok", http.StatusBadRequest, user.ErrorAlreadyExist},
		{"InternalServerError", "create_user_ok", http.StatusInternalServerError, errors.New("fail")},
	}

	for _, tc := range tt {
		// When
		service := &serviceMock{}

		service.On("CreateUser").Return(int64(1), tc.Error)

		rr := httptest.NewRecorder()
		router := gin.Default()
		API(router, service)
		body, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", tc.Filename))
		require.NoError(t, err)
		reader := bytes.NewReader(body)
		request, err := http.NewRequest(http.MethodPost, "/users", reader)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		// Then
		require.Equal(t, tc.ExpectedStatus, rr.Code, "%s failed. Response: %v", tc.TestName, rr.Code)
	}
}

func Test_Handler_API_getUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	userResponse := user.User{
		ID:        1,
		FirstName: "maria",
		LastName:  "garcia",
		Alias:     "mariagarcia",
		Email:     "@gmail",
	}

	tt := []struct {
		TestName       string
		ExpectedStatus int
		Error          error
		UserReponse    user.User
	}{
		{"Ok", http.StatusOK, nil, userResponse},
		{"ErrorUserNotFound", http.StatusNotFound, user.ErrorUserNotFound, user.User{}},
		{"InternalServerError", http.StatusInternalServerError, errors.New("fail"), user.User{}},
	}

	for _, tc := range tt {
		// When
		service := &serviceMock{}

		service.On("GetUser").Return(tc.UserReponse, tc.Error)

		rr := httptest.NewRecorder()
		router := gin.Default()
		API(router, service)

		request, err := http.NewRequest(http.MethodGet, "/users/1", nil)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		// Then
		require.Equal(t, tc.ExpectedStatus, rr.Code, "%s failed. Response: %v", tc.TestName, rr.Code)
	}
}
func Test_Handler_API_createMovement(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tt := []struct {
		TestName, Filename string
		ExpectedStatus     int
		Error              error
	}{
		{"Ok", "create_movement_ok", http.StatusCreated, nil},
		{"WrongFormat", "create_movement_wrong_format", http.StatusBadRequest, nil},
		{"ErrorWrongCurrency", "create_movement_ok", http.StatusBadRequest, movement.ErrorWrongCurrency},
		{"ErrorInsufficientBalance", "create_movement_ok", http.StatusBadRequest, movement.ErrorInsufficientBalance},
		{"InternalServerError", "create_movement_ok", http.StatusInternalServerError, errors.New("fail")},
	}

	for _, tc := range tt {
		// When
		service := &serviceMock{}

		service.On("CreateMovement").Return(int64(1), tc.Error)

		rr := httptest.NewRecorder()
		router := gin.Default()
		API(router, service)
		body, err := ioutil.ReadFile(fmt.Sprintf("testdata/%s.json", tc.Filename))
		require.NoError(t, err)
		reader := bytes.NewReader(body)
		request, err := http.NewRequest(http.MethodPost, "/movements", reader)
		assert.NoError(t, err)

		router.ServeHTTP(rr, request)
		// Then
		require.Equal(t, tc.ExpectedStatus, rr.Code, "%s failed. Response: %v", tc.TestName, rr.Code)
	}
}

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) CreateUser(ctx context.Context, name, lastName, alias, email string) (int64, error) {
	args := s.Called()
	return args.Get(0).(int64), args.Error(1)
}

func (s *serviceMock) GetUser(ctx context.Context, id int64) (user.User, error) {
	args := s.Called()
	return args.Get(0).(user.User), args.Error(1)
}

func (s *serviceMock) CreateMovement(ctx context.Context, movement movement.Movement) (int64, error) {
	args := s.Called()
	return args.Get(0).(int64), args.Error(1)
}
func (s *serviceMock) SearchMovement(ctx context.Context, userID int64, limit, offset uint64, movType, currencyName string) ([]movement.Row, error) {
	args := s.Called()
	return args.Get(0).([]movement.Row), args.Error(1)
}
