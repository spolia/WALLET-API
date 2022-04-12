package movement

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func New(db *sql.DB) *repository {
	return &repository{db: db}
}

// Save inserts a new movement in the user account
func (r repository) Save(ctx context.Context, movement Movement) error {
	var table string
	if table = getCurrencyTable(movement.CurrencyName); table == "" {
		return ErrorWrongCurrency
	}

	query := fmt.Sprintf("INSERT INTO %s(mov_type,currency_name,tx_amount,alias,interaction_alias)VALUES (?,?,?,?,?);", table)

	tx, err := r.db.Begin()
	// sender
	_, err = tx.ExecContext(ctx, query, movement.Type, movement.CurrencyName, movement.Amount, movement.Alias, movement.InteractionAlias)
	if err != nil {
		tx.Rollback()
		return err
	}

	if movement.Alias != movement.InteractionAlias {
		// destiny
		_, err = tx.ExecContext(ctx, query, ReceiveMov, movement.CurrencyName, movement.Amount, movement.InteractionAlias, movement.Alias)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetFunds returns the user funds for a currency
func (r repository) GetFunds(ctx context.Context, currencyName, alias string) (float64, error) {
	var table string
	if table = getCurrencyTable(currencyName); table == "" {
		return 0, ErrorWrongCurrency
	}

	query := fmt.Sprintf("SELECT total_amount FROM %s WHERE id = (SELECT MAX(id) FROM %s WHERE alias = ?);", table, table)
	row := r.db.QueryRowContext(ctx, query, alias)
	var queryResult struct {
		TotalAmount float64
	}

	if err := row.Scan(&queryResult.TotalAmount); err != nil {
		return 0, err
	}

	return queryResult.TotalAmount, nil
}

// InitSave saves initials account movements for a new user
func (r repository) InitSave(ctx context.Context, movement Movement) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	for _, v := range movementTables {
		query := fmt.Sprintf("INSERT INTO %s(mov_type,tx_amount,total_amount,alias,interaction_alias)VALUES (?,?,?,?,?);", v)
		if _, err = tx.ExecContext(ctx, query, movement.Type, movement.Amount, movement.TotalAmount, movement.Alias, movement.InteractionAlias); err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

// GetAccountExtract given an alias returns the funds for all user currencies
func (r repository) GetAccountExtract(ctx context.Context, alias string) (AccountBalance, error) {
	var accountBalance = make(AccountBalance, 0)
	for k, v := range movementTables {
		var queryResult struct {
			TotalAmount float64
		}

		row := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT total_amount FROM %s WHERE id = (SELECT MAX(id) "+
			"FROM %s WHERE alias = ?)", v, v), alias)
		if err := row.Scan(&queryResult.TotalAmount); err != nil {
			return AccountBalance{}, err
		}

		accountBalance[k] = queryResult.TotalAmount
	}

	return accountBalance, nil
}

// GetHistory returns the account history for all the user currencies
func (r repository) GetHistory(ctx context.Context, alias string) (AccountHistory, error) {
	var history = make(AccountHistory, 0)
	for k, v := range movementTables {
		rows, err := r.db.QueryContext(ctx, fmt.Sprintf("SELECT mov_type,date_created,tx_amount,total_amount,interaction_alias "+
			"FROM %s WHERE alias = ?", v), alias)
		if err != nil {
			return AccountHistory{}, err
		}

		var mov []Row
		for rows.Next() {
			var r Row
			err = rows.Scan(&r.Type, &r.DateCreated, &r.Amount, &r.TotalAmount, &r.InteractionAlias)
			if err != nil {
				return AccountHistory{}, err
			}

			mov = append(mov, r)
		}

		history[k] = mov
	}

	return history, nil
}
