package user

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestSave_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	input := User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
		Password:  "1234",
	}
	// When
	mock.ExpectExec("INSERT INTO users(alias,first_name,last_name,email,password) VALUES(?,?,?,?,?);").
		WithArgs(input.Alias, input.FirstName, input.LastName, input.Email, input.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// then
	err = repository.Save(context.Background(), input)
	require.NoError(t, err)
}

func TestSave_Fail(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	input := User{
		FirstName: "name",
		LastName:  "lastname",
		Alias:     "alias",
		Email:     "email",
		Password:  "1234",
	}
	// When
	mock.ExpectExec("INSERT INTO users(alias,first_name,last_name,email,password) VALUES(?,?,?,?,?);").
		WithArgs(input.Alias, input.FirstName, input.LastName, input.Email, input.Password).WillReturnError(&mysql.MySQLError{
		Number: 1062,
	})

	// then
	err = repository.Save(context.Background(), input)
	require.Error(t, err)
	require.EqualError(t, ErrorAlreadyExist, err.Error())
}

func TestDelete_Ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	// When
	mock.ExpectExec("DELETE FROM users Where alias = ?;").
		WithArgs("user").WillReturnResult(sqlmock.NewResult(1, 1))

	// then
	err = repository.Delete(context.Background(), "user")
	require.NoError(t, err)
}

func TestDelete_Fail(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	// When
	mock.ExpectExec("DELETE FROM users Where alias = ?;").
		WithArgs("user").WillReturnError(errors.New("database error"))

	// then
	err = repository.Delete(context.Background(), "user")
	require.Error(t, err)
}

func TestExist_Ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	// When
	mock.ExpectQuery("SELECT alias FROM users WHERE alias = ?;").
		WithArgs("user").WillReturnRows(sqlmock.NewRows([]string{"alias"}).
		AddRow("user"))

	// then
	exist, err := repository.Exist(context.Background(), "user")
	require.NoError(t, err)
	require.True(t, exist)
}
