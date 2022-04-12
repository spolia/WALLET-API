package internal

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

type Service interface {
	CreateUser(ctx context.Context, u user.User) error
	GetBalance(ctx context.Context, alias string) (movement.AccountBalance, error)
	Send(ctx context.Context, m movement.Movement) error
	AutoDeposit(ctx context.Context, m movement.Movement) error
	GetHistory(ctx context.Context, alias string) (movement.AccountHistory, error)
	ValidateCredential(ctx context.Context, alias, password string) (bool, error)
}

func API(r *mux.Router, service Service) {
	r.HandleFunc("/login", login(service)).Methods(http.MethodPost)
	r.HandleFunc("/logout", logout).Methods(http.MethodPost)
	r.HandleFunc("/internal/movements/balance", getBalance(service)).Methods(http.MethodGet)
	r.HandleFunc("/internal/movements/history", getHistory(service)).Methods(http.MethodGet)
	r.HandleFunc("/internal/movements/send", send(service)).Methods(http.MethodPost)

	// Useful to test with more users and fund accounts of different currencies
	r.HandleFunc("/users", createUser(service)).Methods(http.MethodPost)
	r.HandleFunc("/internal/movements/deposit", deposit(service)).Methods(http.MethodPost)
}

// hacer la respuesta de history mas linda
// encrypt passport
