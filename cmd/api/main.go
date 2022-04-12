package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spolia/wallet-api/cmd/api/internal"
	"github.com/spolia/wallet-api/internal/wallet"
	"github.com/spolia/wallet-api/internal/wallet/movement"
	"github.com/spolia/wallet-api/internal/wallet/user"
)

func main() {
	log.Println("starting")

	db, err := sql.Open("mysql", "tester:secret@tcp(db:3306)/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	service := wallet.New(user.New(db), movement.New(db))
	log.Println("service successfully configured")

	router := mux.NewRouter()
	internal.API(router, service)
	// localhost:8080
	http.ListenAndServe(":8080", router)
}
