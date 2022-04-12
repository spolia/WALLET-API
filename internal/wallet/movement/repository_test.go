package movement

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestSaveMovement_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()
	movement := Movement{
		Type:             "send",
		Amount:           100.2,
		CurrencyName:     USDT,
		Alias:            "user",
		InteractionAlias: "otheruser",
	}

	// When
	mock.ExpectBegin()
	query1 := "INSERT INTO movements_usdt(mov_type,currency_name,tx_amount,alias,interaction_alias)VALUES (?,?,?,?,?);"
	mock.ExpectExec(query1).WithArgs(movement.Type, movement.CurrencyName, movement.Amount, movement.Alias, movement.InteractionAlias).WillReturnResult(sqlmock.NewResult(1, 1))
	query2 := "INSERT INTO movements_usdt(mov_type,currency_name,tx_amount,alias,interaction_alias)VALUES (?,?,?,?,?);"
	mock.ExpectExec(query2).WithArgs("receive", movement.CurrencyName, movement.Amount, movement.InteractionAlias, movement.Alias).WillReturnResult(sqlmock.NewResult(2, 1))
	mock.ExpectCommit()
	// then
	err = repository.Save(context.Background(), movement)
	require.NoError(t, err)
}

func TestSaveMovement_Error(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:         "send",
		Amount:       100.2,
		CurrencyName: USDT,
		Alias:        "user",
	}

	// When
	mock.ExpectBegin()
	query1 := "INSERT INTO movements_usdt(mov_type,currency_name,tx_amount,alias,interaction_alias)VALUES (?,?,?,?,?);"
	mock.ExpectExec(query1).WithArgs(movement.Type, movement.CurrencyName, movement.Amount, movement.Alias, movement.InteractionAlias).WillReturnError(&mysql.MySQLError{
		Number: 1264,
	})
	mock.ExpectRollback()
	// then
	err = repository.Save(context.Background(), movement)
	require.Error(t, err)
}

func TestSaveMovement_ErrorWrongCurrency(t *testing.T) {
	// Given
	repository := New(nil)

	// When
	err := repository.Save(context.Background(), Movement{
		Type:         DepositMov,
		Amount:       100.2,
		CurrencyName: "wrong",
		Alias:        "alias",
	})

	// Then
	require.Error(t, err)
	require.EqualError(t, ErrorWrongCurrency, err.Error())
}
func TestGetFunds_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()
	movement := Movement{
		Type:             "send",
		Amount:           100.2,
		CurrencyName:     USDT,
		Alias:            "user",
		InteractionAlias: "user",
	}

	// When
	query := "SELECT total_amount FROM movements_usdt WHERE id = (SELECT MAX(id) FROM movements_usdt WHERE alias = ?);"
	mock.ExpectQuery(query).WithArgs(movement.Alias).WillReturnRows(sqlmock.NewRows([]string{"total_amount"}).
		AddRow(float64(100)))
	// then
	result, err := repository.GetFunds(context.Background(), movement.CurrencyName, movement.Alias)
	require.NoError(t, err)
	require.Equal(t, float64(100), result)
}
