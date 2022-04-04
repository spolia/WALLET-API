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
	}
	// When
	mock.ExpectExec("INSERT INTO users(first_name,last_name,alias,email)VALUES (?,?,?,?);").
		WithArgs(input.FirstName, input.LastName, input.Alias, input.Email).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// then
	id, err := repository.Save(context.Background(), input.FirstName, input.LastName, input.Alias, input.Email)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)
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
	}
	// When
	mock.ExpectExec("INSERT INTO users(first_name,last_name,alias,email)VALUES (?,?,?,?);").
		WithArgs(input.FirstName, input.LastName, input.Alias, input.Email).WillReturnError(&mysql.MySQLError{
		Number: 1062,
	})

	// then
	id, err := repository.Save(context.Background(), input.FirstName, input.LastName, input.Alias, input.Email)
	require.Error(t, err)
	require.EqualError(t, ErrorAlreadyExist, err.Error())
	require.Equal(t, int64(0), id)
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
	mock.ExpectExec("DELETE FROM users Where id = ?;").
		WithArgs(int64(1)).WillReturnResult(sqlmock.NewResult(1, 1))

	// then
	err = repository.Delete(context.Background(), 1)
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
	mock.ExpectExec("DELETE FROM users Where id = ?;").
		WithArgs(int64(1)).WillReturnError(errors.New("database error"))

	// then
	err = repository.Delete(context.Background(), 1)
	require.Error(t, err)
}

func TestGet_Ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	// When
	mock.ExpectQuery("SELECT * FROM users Where id = ?;").
		WithArgs(int64(1)).WillReturnRows(sqlmock.NewRows([]string{"id", "first_name", "last_name", "alias", "email"}).
		AddRow(1, "maria", "garcia", "alias", "@gmail"))

	// then
	userResponse, err := repository.Get(context.Background(), int64(1))
	require.NoError(t, err)
	require.NotEmpty(t, userResponse)
}
