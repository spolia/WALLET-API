package movement

import (
	"context"
	"errors"
	"time"
)

const (
	DepositMov = "deposit"
	BTC        = "BTC"
	ARS        = "ARS"
	USDT       = "USDT"
)

var movementTables = map[string]string{
	USDT: "movements_usdt",
	ARS:  "movements_ars",
	BTC:  "movements_btc",
}

var (
	ErrorInsufficientBalance = errors.New("movement: insufficient balance")
	ErrorWrongOperation      = errors.New("movement: wrong operation")
	ErrorWrongUser           = errors.New("movement: wrong user")
	ErrorWrongCurrency       = errors.New("movement: wrong currency")
	ErrorNoMovements         = errors.New("movement: there no movements")
)

type AccountBalance map[string]float64
type AccountHistory map[string][]MovRow

type MovRow struct {
	CurrencyName     string
	Type             string
	DateCreated      time.Time
	Amount           float64
	TotalAmount      float64
	InteractionAlias string
}

type Repository interface {
	Save(ctx context.Context, movement Movement) (int64, error)
	InitSave(ctx context.Context, movement Movement) error
	GetAccountExtract(ctx context.Context, alias string) (AccountBalance, error)
}

type Movement struct {
	ID               int64   `json:"id"`
	Type             string  `json:"type" binding:"required,oneof=deposit extract"`
	Amount           float64 `json:"amount" binding:"required,gte=0"`
	CurrencyName     string  `json:"currencyname" binding:"required,oneof=usdt btc ars"`
	Alias            string  `json:"alias" binding:"required"`
	TotalAmount      float64 `json:"totalamount"`
	InteractionAlias string  `json:"interactionalias" binding:"required"`
}

type Currency struct {
	ID     int64
	Name   string
	Digits float64
}

type Row struct {
	CurrencyName string
	Type         string
	DateCreated  time.Time
	Amount       float64
	TotalAmount  float64
}

func getCurrencyTable(currency string) string {
	return movementTables[currency]
}

func getCurrenciesTables(currency string) []string {
	var tables = make([]string, 0)

	if movementTables[currency] == "" {
		for _, v := range movementTables {
			tables = append(tables, v)
		}
	} else {
		tables = append(tables, movementTables[currency])
	}

	return tables
}
