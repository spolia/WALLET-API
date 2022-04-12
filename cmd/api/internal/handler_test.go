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

	"github.com/gorilla/mux"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Handler_API_createUser(t *testing.T) {
	tt := []struct {
		TestName, Filename string
		ExpectedStatus     int
		Error              error
	}{
		{"Ok", "create_user_ok", http.StatusOK, nil},
		{"WrongFormat", "create_user_wrong_format", http.StatusBadRequest, nil},
		{"ErrorAlreadyExist", "create_user_ok", http.StatusBadRequest, user.ErrorAlreadyExist},
		{"InternalServerError", "create_user_ok", http.StatusInternalServerError, errors.New("fail")},
	}

	for _, tc := range tt {
		// When
		service := &serviceMock{}

		service.On("CreateUser").Return(tc.Error)

		rr := httptest.NewRecorder()
		router := mux.NewRouter()
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

func Test_Handler_API_Login(t *testing.T) {

	service := &serviceMock{}

	service.On("ValidateCredential").Return(true, nil)

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	API(router, service)
	body, err := ioutil.ReadFile("testdata/login.json")
	require.NoError(t, err)
	reader := bytes.NewReader(body)
	request, err := http.NewRequest(http.MethodPost, "/login", reader)
	assert.NoError(t, err)

	router.ServeHTTP(rr, request)
	// Then
	require.Equal(t, http.StatusOK, rr.Code)
}

type serviceMock struct {
	mock.Mock
}

func (s *serviceMock) CreateUser(ctx context.Context, u user.User) error {
	args := s.Called()
	return args.Error(0)
}

func (s *serviceMock) GetBalance(ctx context.Context, alias string) (movement.AccountBalance, error) {
	args := s.Called()
	return args.Get(0).(movement.AccountBalance), args.Error(1)
}

func (s *serviceMock) Send(ctx context.Context, m movement.Movement) error {
	args := s.Called()
	return args.Error(0)
}

func (s *serviceMock) AutoDeposit(ctx context.Context, m movement.Movement) error {
	args := s.Called()
	return args.Error(0)
}

func (s *serviceMock) GetHistory(ctx context.Context, alias string) (movement.AccountHistory, error) {
	args := s.Called()
	return args.Get(0).(movement.AccountHistory), args.Error(1)
}

func (s *serviceMock) ValidateCredential(ctx context.Context, alias, password string) (bool, error) {
	args := s.Called()
	return args.Bool(0), args.Error(1)
}
