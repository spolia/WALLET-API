package movement

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/spolia/wallet-api/internal/database"
	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSaveMovement_ok(t *testing.T) {
	// Given
	// Work out the path to the 'scripts' directory and set mount strings
	packageName := "database"
	workingDir, _ := os.Getwd()
	rootDir := strings.Replace(workingDir, packageName, "", 1)
	mountFrom := fmt.Sprintf("%s/migrations/init.sql", rootDir)
	mountTo := "/docker-entrypoint-initdb.d/init.sql"
	// Create the Postgres TestContainer
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:11.6-alpine",
		ExposedPorts: []string{"5432/tcp"},
		BindMounts:   map[string]string{mountFrom: mountTo},
		Env: map[string]string{
			"POSTGRES_DB": os.Getenv("DBNAME"),
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(err)
	defer postgresC.Terminate(ctx)

	// Get the port mapped to 5432 and set as ENV
	p, _ := postgresC.MappedPort(ctx, "5432")
	os.Setenv("DBPORT", p.Port())

	db, err := database.InitDB()
	require.NoError(err)

	repository := New(db)

	movement := Movement{
		Type:             DepositMov,
		Amount:           100.2,
		CurrencyName:     USDT,
		Alias:            "user",
		InteractionAlias: "user",
	}

	// When
	// then
	movementID, err := repository.Save(ctx, movement)
	require.NoError(t, err)
	require.Equal(t, int64(1), movementID)
}

/*
func TestSaveMovement_Error(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:         DepositMov,
		Amount:       100.2,
		CurrencyName: USDT,
		UserID:       1,
	}

	// When
	query := "INSERT INTO movements_usdt(mov_type,currency_name,tx_amount,user_id)VALUES (?,?,?,?);"
	mock.ExpectExec(query).
		WithArgs(movement.Type, movement.CurrencyName, movement.Amount, movement.UserID).WillReturnError(&mysql.MySQLError{
		Number: 1264,
	})

	// then
	movementID, err := repository.Save(context.Background(), movement)
	require.Error(t, err)
	require.EqualError(t, ErrorInsufficientBalance, err.Error())
	require.Equal(t, int64(0), movementID)
}

func TestSaveMovement_ErrorWrongCurrency(t *testing.T) {
	// Given
	repository := New(nil)

	// When
	movementID, err := repository.Save(context.Background(), Movement{
		Type:         DepositMov,
		Amount:       100.2,
		CurrencyName: "wrong",
		UserID:       1,
	})

	// Then
	require.Error(t, err)
	require.EqualError(t, ErrorWrongCurrency, err.Error())
	require.Equal(t, int64(0), movementID)
}

func TestInitSave_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:   DepositMov,
		UserID: 1,
	}
	// When
	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO movements_usdt(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);").
		WithArgs(movement.Type, movement.Amount, movement.TotalAmount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO movements_ars(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);").
		WithArgs(movement.Type, movement.Amount, movement.TotalAmount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("INSERT INTO movements_btc(mov_type,tx_amount,total_amount,user_id)VALUES (?,?,?,?);").
		WithArgs(movement.Type, movement.Amount, movement.TotalAmount, movement.UserID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// then
	err = repository.InitSave(context.Background(), movement)
	require.NoError(t, err)
}

func TestSearch_ok(t *testing.T) {
	// Given
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		require.NoError(t, err)
	}
	repository := New(db)
	defer db.Close()

	movement := Movement{
		Type:         DepositMov,
		UserID:       int64(1),
		CurrencyName: ARS,
	}
	// When

	mock.ExpectQuery("SELECT mov_type, currency_name, date_created, tx_amount, total_amount FROM movements_ars WHERE user_id = ? AND mov_type = 'deposit'").
		WithArgs(movement.UserID).WillReturnRows(sqlmock.NewRows([]string{"mov_type", "currency_name", "date_created", "tx_amount", "total_amount"}).
		AddRow("deposit", "ars", time.Now(), 200, 1000)).WillReturnRows(sqlmock.NewRows([]string{"mov_type", "currency_name", "date_created", "tx_amount", "total_amount"}).
		AddRow("deposit", "ars", time.Now(), 300, 2000))

	// then
	rows, err := repository.Search(context.Background(), movement.UserID, 0, 0, movement.Type, movement.CurrencyName)
	require.NoError(t, err)
	require.True(t, true, len(rows) > 0)
}
*/
