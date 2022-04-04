package movement

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

// Save inserts a new movement in the database
func (r repository) Save(ctx context.Context, movement Movement) (int64, error) {
	var table string
	if table = getCurrencyTable(movement.CurrencyName); table == "" {
		return 0, ErrorWrongCurrency
	}

	query := fmt.Sprintf("INSERT INTO %s(mov_type,currency_name,tx_amount,alias,interaction_alias)VALUES ($1,$2,$3,$4,$5);", table)

	result, err := r.db.ExecContext(ctx, query, movement.Type, movement.CurrencyName, movement.Amount, movement.Alias, movement.InteractionAlias)
	if err != nil {
		// when tx_amount - total_amount is less than 0
		if err.(*mysql.MySQLError).Number == 1264 || err.(*mysql.MySQLError).Number == 1690 {
			return 0, ErrorInsufficientBalance
		}
		// wrong type
		if err.(*mysql.MySQLError).Number == 1265 {
			return 0, ErrorWrongOperation
		}

		if err.(*mysql.MySQLError).Number == 1048 {
			return 0, ErrorWrongUser
		}

		return 0, err
	}
	id, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InitSave saves initials movements for a new user
func (r repository) InitSave(ctx context.Context, movement Movement) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	for _, v := range movementTables {
		query := fmt.Sprintf("INSERT INTO %s(mov_type,tx_amount,total_amount,alias,interaction_alias)VALUES ($1,$2,$3,$4,$5);", v)

		if _, err = tx.ExecContext(ctx, query, nil, movement.Type, movement.Amount, movement.TotalAmount, movement.Alias, movement.InteractionAlias); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAccountExtract given an alias returns the last movements for each currency
func (r repository) GetAccountExtract(ctx context.Context, alias string) (AccountBalance, error) {
	var accountBalance = make(AccountBalance, 0)
	for k, v := range movementTables {
		var queryResult struct {
			totalAmount float64
		}

		row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT total_amount FROM %s WHERE date_created = (SELECT MAX(date_created) "+
			"FROM %s WHERE alias = $1)", v, v), nil, alias)
		if err := row.Scan(&queryResult.totalAmount); err != nil {
			return AccountBalance{}, err
		}

		accountBalance[k] = queryResult.totalAmount
	}

	return accountBalance, nil
}

func (r repository) GetHistory(ctx context.Context, alias string) (AccountHistory, error) {
	var history = make(AccountHistory, 0)
	for k, v := range movementTables {
		var result = make([]MovRow, 0)

		rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT mov_type,date_created,tx_amount,total_amount,interaction_alias "+
			"FROM %s WHERE alias = $1", v), nil, alias)
		if err != nil {
			return AccountHistory{}, err
		}

		if err := rows.Scan(pq.Array(&result)); err != nil {
			return AccountHistory{}, err
		}

		history[k] = result
	}

	return history, nil
}
