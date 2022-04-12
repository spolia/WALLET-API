package user

import (
	"context"
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

// Save inserts a new user
func (r repository) Save(ctx context.Context, u User) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users(alias,first_name,last_name,email,password) VALUES(?,?,?,?,?);",
		u.Alias, u.FirstName, u.LastName, u.Email, u.Password)
	if err != nil {
		if v, ok := err.(*mysql.MySQLError); ok {
			if v.Number == 1062 {
				return ErrorAlreadyExist
			}
		}

		return err
	}

	return nil
}

// Delete deletest an user
func (r repository) Delete(ctx context.Context, alias string) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users Where alias = ?;", alias)
	if err != nil {
		return err
	}

	if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

// Exist returns true if a user with the given alias exist in the database
func (r repository) Exist(ctx context.Context, alias string) (bool, error) {
	row := r.db.QueryRowContext(ctx, "SELECT alias FROM users WHERE alias = ?;", alias)
	if row.Err() != nil {
		return false, row.Err()
	}

	var queryResult struct {
		Alias string
	}

	if err := row.Scan(&queryResult.Alias); err != nil {
		if err == sql.ErrNoRows {
			return false, ErrorDestinyUserNotFound
		}
		return false, err
	}

	return queryResult.Alias != "", nil
}

// IsValidCredential returns true a user with the given alias and password exist in the database
func (r repository) IsValidCredential(ctx context.Context, alias, password string) (bool, error) {
	row := r.db.QueryRowContext(ctx, "SELECT alias FROM users WHERE alias = ? AND password = ?;", alias, password)
	if row.Err() != nil {
		return false, row.Err()
	}

	var queryResult struct {
		Alias string
	}

	if err := row.Scan(&queryResult.Alias); err != nil {
		if err == sql.ErrNoRows {
			return false, ErrorInvalidCredential
		}

		return false, err
	}

	return queryResult.Alias != "", nil
}
