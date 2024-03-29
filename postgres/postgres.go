package postgres

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Client enables a connection to a postgres db
type Client struct {
	db         *gorm.DB
	connString string
}

// NewClient starts a connection with db using gorm.
func NewClient(host string, port int, user, password, dbname string, ssl bool) (*Client, error) {
	// Connection string
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require sslrootcert=/etc/ssl/certs/ca.crt",
		host, port, user, password, dbname)

	if !ssl {
		connString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
	}

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
		db: db.Debug(),

		connString: connString,
	}

	return client, nil
}

// Repositories contains all repositories of the application
type Repositories struct {
	UserPosition *UserPositionRepository
}

// NewRepositories creates new repositories
func NewRepositories(client *Client) (*Repositories, error) {
	repos := &Repositories{
		UserPosition: NewUserPositionRepository(client),
	}

	return repos, nil
}

// MigrateUp runs migrations
func MigrateUp(host string, port int, user, password, dbname string, ssl bool	) error {
	// Connection string
	connString := fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=require&sslrootcert=\/etc\/ssl\/certs\/ca.crt`,
		user, url.QueryEscape(password), host, port, dbname)

	if !ssl {
		connString = fmt.Sprintf(`postgres://%s:%s@%s:%d/%s?sslmode=disable`,
			user, url.QueryEscape(password), host, port, dbname)
	}

	u, _ := url.Parse(connString)
	fmt.Println("ssl is", u.Query().Get("sslrootcert"))
	dbm := dbmate.New(u)

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	currentFolder := filepath.Base(dir)

	configPath := dir + "/db/migrations"
	if currentFolder == "nawaltni" {
		configPath = dir + "/tracker/db/migrations"
	}

	dbm.MigrationsDir = []string{configPath}
	err = dbm.CreateAndMigrate()
	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Printf("Migrations ran successfully")

	return nil
}
