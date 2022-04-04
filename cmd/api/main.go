package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spolia/wallet-api/cmd/api/internal"
	"github.com/spolia/wallet-api/internal/wallet"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

func main() {
	log.Println("starting")
	var err error
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", "root", "rootroot", "127.0.0.1:3306", "wallet", "parseTime=true")
	log.Println("datasource", dataSourceName)

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	service := wallet.New(user.New(db), movement.New(db))
	log.Println("service successfully configured")

	router := gin.Default()
	internal.API(router, service)

	router.Run("localhost:8080")
	log.Println("listening")
}
