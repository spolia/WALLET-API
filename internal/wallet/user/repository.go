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
func (r repository) Save(ctx context.Context, firstName, lastName, alias, email string) (int64, error) {
	result, err := r.db.ExecContext(ctx, "INSERT INTO users(first_name,last_name,alias,email)VALUES (?,?,?,?);",
		firstName, lastName, alias, email)
	if err != nil {
		if err.(*mysql.MySQLError).Number == 1062 {
			return 0, ErrorAlreadyExist
		}
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// Delete deletest an user
func (r repository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, "DELETE FROM users Where id = ?;", id)
	if err != nil {
		return err
	}

	if _, err = result.RowsAffected(); err != nil {
		return err
	}

	return nil
}

// Get returns a user
func (r repository) Get(ctx context.Context, id int64) (User, error) {
	row := r.db.QueryRowContext(ctx, "SELECT * FROM users Where id = ?;", id)
	if row.Err() != nil {
		return User{}, row.Err()
	}

	var user User
	if err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Alias, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			return User{}, ErrorUserNotFound
		}
		return User{}, err
	}

	return user, nil
}
