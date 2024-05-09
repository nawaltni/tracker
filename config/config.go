package config

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Config represents the application config
type Config struct {
	Environment string `mapstructure:"environment"`
	// Logging configurations
	LogConfig Logger `mapstructure:"log_config"`
	Postgres  struct {
		Host          string `mapstructure:"host"`
		Port          int    `mapstructure:"port"`
		Username      string `mapstructure:"user"`
		Password      string `mapstructure:"password"`
		Database      string `mapstructure:"database"`
		RunMigrations bool   `mapstructure:"run_migrations"`
		SSL           bool   `mapstructure:"ssl"`
	}

	HTTP struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	GRPC struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"grpc"`

	Places   Places   `mapstructure:"places"`
	Auth     Auth     `mapstructure:"auth"`
	Bigquery Bigquery `mapstructure:"bigquery"`
}

type Places struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Auth struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Bigquery struct {
	ProjectID string `mapstructure:"project_id"`
	DataSetID string `mapstructure:"dataset_id"`
}

// LoadConfig loads the configuration from file and environment variables
func LoadConfig(cmd *cobra.Command) (*Config, error) {
	// from the command itself
	if err := viper.BindPFlags(cmd.Flags()); err != nil {
		return nil, err
	}

	// from the environment
	viper.SetEnvPrefix("TRACKER")
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

	config := new(Config)
	if err := viper.Unmarshal(config); err != nil {
		return nil, err
	}

	// all the error signatures above had to change to nil, err
	return config, nil
}
