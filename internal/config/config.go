package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the configuration structure
type AuthConfig struct {
	Username     string `mapstructure:"username"`
	Token        string `mapstructure:"token"`
	ProfileID    string `mapstructure:"profile_id"`
	RefreshToken string `mapstructure:"refresh_token"`
	EndPoint     string `mapstructure:"endpoint"`
}

type Config interface {
	UpdateEndpoint(string) error
	UpdateConfig(string, string, string, string) error
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
	viper.SetDefault("endpoint", "https://api.logfire.ai/")

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
func (c *cfg) UpdateConfig(username, token, profileID string, refreshToken string) error {
	// Write the updated configuration to the file
	viper.Set("username", username)
	viper.Set("token", token)
	viper.Set("profile_id", profileID)
	viper.Set("refresh_token", refreshToken)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	c.AuthCfg.Username = username
	c.AuthCfg.Token = token
	c.AuthCfg.ProfileID = profileID
	c.AuthCfg.RefreshToken = refreshToken
	return nil
}

func (c *cfg) UpdateEndpoint(endpoint string) error {
	viper.Set("endpoint", endpoint)

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	c.AuthCfg.EndPoint = endpoint
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
