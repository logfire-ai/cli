package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the configuration structure
type AuthConfig struct {
	Username  string `mapstructure:"username"`
	Token     string `mapstructure:"token"`
	ProfileID string `mapstructure:"profile_id"`
}

type Config interface {
	UpdateConfig(string, string, string) error
	DeleteConfig() error
	HasEnvToken() bool
	Get() *AuthConfig
}

type cfg struct {
	AuthCfg *AuthConfig
}

func NewConfig() (Config, error) {
	configFile, err := getHomeConfigPath()
	if err != nil {
		return nil, err
	}

	// Set up Viper for YAML configuration
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	viper.SetDefault("username", "")
	viper.SetDefault("token", "")
	viper.SetDefault("profile_id", "")

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

	return &cfg{authConfig}, nil
}

func (c *cfg) Get() *AuthConfig {
	return c.AuthCfg
}

// UpdateConfig updates the configuration values and writes them to the config file
func (c *cfg) UpdateConfig(username, token, profileID string) error {
	// Write the updated configuration to the file
	viper.Set("username", username)
	viper.Set("token", token)
	viper.Set("profile_id", profileID)

	// Write the updated configuration to the file
	if err := viper.WriteConfig(); err != nil {
		return err
	}

	c.AuthCfg.Username = username
	c.AuthCfg.Token = token
	c.AuthCfg.ProfileID = profileID
	return nil
}

func (c *cfg) HasEnvToken() bool {
	return c.AuthCfg.Token != "" || c.AuthCfg.Username != "" || c.AuthCfg.ProfileID != ""
}

func (c *cfg) DeleteConfig() error {
	c.AuthCfg.Username = ""
	c.AuthCfg.Token = ""
	c.AuthCfg.ProfileID = ""

	configFile, err := getHomeConfigPath()
	if err != nil {
		return err
	}

	// Delete file from homedir.
	err = os.Remove(configFile)
	if err != nil {
		return err
	}

	return nil
}

func getHomeConfigPath() (string, error) {
	// Get user's home directory
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// Add .logfire to file path
	configFile := filepath.Join(usr.HomeDir, ".logfire")
	return configFile, nil
}
