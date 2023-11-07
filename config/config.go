package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Configuration represents the application configuration
type Configuration struct {
	Database struct {
		Username string
		Password string
		DBName   string
		Host     string
		Port     int
	}
	GRPC struct {
		Host string
		Port int
	}
}

// LoadConfig loads the configuration from file and environment variables
func LoadConfig(cmd *cobra.Command) (*Configuration, error) {
	// from the command itself
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	// from the environment
	viper.SetEnvPrefix("PLACES")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	configFile, err := cmd.Flags().GetString("config")
	if err != nil {
		return nil, err
	}

	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// from a config file
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	// NOTE: this will require that you have config file somewhere in the paths specified. It can be reading from JSON, TOML, YAML, HCL, and Java properties files.
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	config := new(Configuration)
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// all the error signatures above had to change to nil, err
	return config, nil
}
