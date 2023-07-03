package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the configuration structure
type AuthConfig struct {
	Username string `mapstructure:"username"`
	Token    string `mapstructure:"token"`
}

type Config struct {
	AuthConfig *AuthConfig
}

// InitializeConfig initializes the configuration file
func NewConfig() (*Config, error) {
	// Get user's home directory
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	configFile := filepath.Join(usr.HomeDir, ".logfire.yaml")

	// Set up Viper for YAML configuration
	viper.SetConfigFile(configFile)
	viper.SetDefault("username", "")
	viper.SetDefault("token", "")

	// Check if the config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Config file doesn't exist, create a new one
		if err := viper.WriteConfig(); err != nil {
			return nil, err
		}
	}

	// Read the configuration from the file
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// Unmarshal the configuration into the Config struct
	authConfig := &AuthConfig{}
	if err := viper.Unmarshal(authConfig); err != nil {
		return nil, err
	}

	return &Config{AuthConfig: authConfig}, nil
}

// UpdateConfig updates the configuration values and writes them to the config file
func (c *Config) UpdateConfig(username, token string) error {
	// Write the updated configuration to the file
	viper.Set("username", username)
	viper.Set("token", token)

	// Write the updated configuration to the file
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	c.AuthConfig.Username = username
	c.AuthConfig.Token = token
	return nil
}
