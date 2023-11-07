package postgres

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Client enables a connection to a postgres db
type Client struct {
	db         *gorm.DB
	connString string
}

// NewClient starts a connection with db using gorm.
func NewClient(user, password, host string, port int, dbname string) (*Client, error) {
	// Connection string
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(connString))
	if err != nil {
		return nil, err
	}

	database, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error retrieving db from gorm: %w", err)
	}

	// Wait for db is up. Try to ping it 10 times.
	for i := 1; ; i++ {
		err = database.Ping()
		if err != nil {
			if i <= 10 {
				time.Sleep(time.Second * 1)
				continue
			}
			return nil, fmt.Errorf("could not establish connection to pg after %d tries", i)
		}
		break
	}
	// Show message of connection when it works
	fmt.Println("Connection successfully")

	client := &Client{
		db: db,

		connString: connString,
	}

	return client, nil
}
