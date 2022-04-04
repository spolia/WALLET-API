package database

import (
	"database/sql"
	"fmt"
	"log"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

// Config is a struct that pulls in env vars to configure the database
type Config struct {
	User string `envconfig:"DBUSER"`
	Name string `envconfig:"DBNAME"`
	Host string `envconfig:"DBHOST"`
	Port string `envconfig:"DBPORT"`
}

// InitDB connects to the database
func InitDB() (*sql.DB, error) {
	var c Config
	err := envconfig.Process("myapp", &c)
	if err != nil {
		return nil, err
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Name)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = backoff.Retry(func() error {
		err := db.Ping()
		if err != nil {
			log.Println("DB is not ready...backing off...")
			return err
		}
		log.Println("DB is ready!")
		return nil
	}, backoff.NewExponentialBackOff())

	if err != nil {
		return nil, err
	}

	return db, nil
}
