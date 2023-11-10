package cmd

import (
	"fmt"
	"log"

	"github.com/nawaltni/tracker/config"
	"github.com/nawaltni/tracker/grpc"
	"github.com/nawaltni/tracker/postgres"
	"github.com/nawaltni/tracker/services"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Tracker is the service that manages tracker data",
	Run:   run,
}

func RootCommand() *cobra.Command {
	// main flags
	rootCmd.Flags().IntP("port", "p", 8080, "Set the port to run the app")
	rootCmd.Flags().StringP("config", "c", "config.toml", "Set the config file to use")
	return rootCmd
}

func run(cmd *cobra.Command, args []string) {
	// 1. Read Config
	conf, err := config.LoadConfig(cmd)
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
		return
	}

	pConfig := conf.Postgres
	fmt.Printf("config: %+v\n", conf)

	// 2. Create Postgres Client
	pgClient, err := postgres.NewClient(pConfig.Host, pConfig.Port, pConfig.Username, pConfig.Password, pConfig.Database)
	if err != nil {
		log.Fatal("Failed to create postgres client: " + err.Error())
		return
	}

	// 3. Migration
	if pConfig.RunMigrations {
		err = postgres.MigrateUp(pConfig.Host, pConfig.Port, pConfig.Username, pConfig.Password, pConfig.Database)
		if err != nil {
			log.Fatal("Failed to run migrations: " + err.Error())
			return
		}
	}

	// 3. Create Repositories
	repost, err := postgres.NewRepositories(pgClient)
	if err != nil {
		log.Fatal("Failed to create postgres repositories: " + err.Error())
		return

	}

	// 4. Create Services
	services, err := services.NewServices(*conf, repost)
	if err != nil {
		log.Fatal("Failed to create services: " + err.Error())
		return
	}

	// 5. Start gRPC Service
	server, err := grpc.New(*conf, services)
	if err != nil {
		log.Fatal("Failed to create grpc server: " + err.Error())
		return
	}

	// 19. Start gRPC Service
	err = server.Start()
	if err != nil {
		log.Fatal(errors.Wrap(err, "Could not start grpc server"))
	}
}
