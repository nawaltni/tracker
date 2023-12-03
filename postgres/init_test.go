package postgres

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/amacneil/dbmate/v2/pkg/dbmate"
	_ "github.com/amacneil/dbmate/v2/pkg/driver/postgres"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	fmt.Println("Starting postgres test container...")
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls an image, creates a container based on it, and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgis/postgis",
		Tag:        "latest",
		Env:        []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres", "POSTGRES_DB=tracker"},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://postgres:postgres@%s/tracker?sslmode=disable", hostAndPort)

	log.Println("Connecting to database on url: ", databaseUrl)

	err = resource.Expire(60) // Tell docker to hard kill the container in 60 seconds
	if err != nil {
		log.Fatalf("Could not set container expiration: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 60 * time.Second

	// exponential backoff-retry for connecting to the database
	if err := pool.Retry(func() error {
		var err error
		db, err = gorm.Open(postgres.Open(databaseUrl), &gorm.Config{})
		if err != nil {
			return err
		}
		db = db.Debug()
		return db.Exec("SELECT 1").Error
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Migrate the schema
	u, _ := url.Parse(databaseUrl)
	dbm := dbmate.New(u)
	dbm.MigrationsDir = []string{"../db/migrations"}
	err = dbm.CreateAndMigrate()
	if err != nil {
		panic(fmt.Errorf("failed to migrate database: %w", err))
	}

	// Run the tests
	code := m.Run()

	// You can force cleanup after tests, but be aware that it's not instant
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
