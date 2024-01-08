package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/viper"
)

// AuthConfig Config represents the configuration structure
type AuthConfig struct {
	Username      string `mapstructure:"username"`
	Role          string `mapstructure:"role"`
	Token         string `mapstructure:"token"`
	ProfileID     string `mapstructure:"profile_id"`
	RefreshToken  string `mapstructure:"refresh_token"`
	EndPoint      string `mapstructure:"endpoint"`
	TeamId        string `mapstructure:"team_id"`
	AccountId     string `mapstructure:"account_id"`
	GrpcEndpoint  string `mapstructure:"grpc_endpoint"`
	GrpcIngestion string `mapstructure:"grpc_ingestion"`
	Theme         string `mapstructure:"theme"`
}

type Config interface {
	UpdateConfig(*string, *string, *string, *string, *string, *string, *string, *string, *string, *string, *string) error
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
	viper.SetDefault("role", "")
	viper.SetDefault("token", "")
	viper.SetDefault("profile_id", "")
	viper.SetDefault("endpoint", "https://api.logfire.ai/")
	viper.SetDefault("grpc_endpoint", "api.logfire.ai:443")
	viper.SetDefault("grpc_ingestion", "https://in.logfire.ai")
	viper.SetDefault("team_id", "")
	viper.SetDefault("account_id", "")
	viper.SetDefault("theme", "dark")

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
func (c *cfg) UpdateConfig(username, role, token, profileID, refreshToken, teamid, accountid, endpoint, grpcEndpoint, grpcIngestion, theme *string) error {
	// Write the updated configuration to the file

	if username != nil {
		viper.Set("username", *username)
		c.AuthCfg.Username = *username
	}

	if role != nil {
		viper.Set("role", *role)
		c.AuthCfg.Role = *role
	}

	if token != nil {
		viper.Set("token", *token)
		c.AuthCfg.Token = *token
	}

	if profileID != nil {
		viper.Set("profile_id", *profileID)
		c.AuthCfg.ProfileID = *profileID
	}

	if refreshToken != nil {
		viper.Set("refresh_token", *refreshToken)
		c.AuthCfg.RefreshToken = *refreshToken
	}

	if teamid != nil {
		viper.Set("team_id", *teamid)
		c.AuthCfg.TeamId = *teamid
	}

	if accountid != nil {
		viper.Set("account_id", *accountid)
		c.AuthCfg.AccountId = *accountid
	}

	if endpoint != nil {
		viper.Set("endpoint", *endpoint)
		c.AuthCfg.EndPoint = *endpoint
	}

	if grpcEndpoint != nil {
		viper.Set("grpc_endpoint", *grpcEndpoint)
		c.AuthCfg.GrpcEndpoint = *grpcEndpoint
	}

	if grpcIngestion != nil {
		viper.Set("grpc_ingestion", *grpcIngestion)
		c.AuthCfg.GrpcIngestion = *grpcIngestion
	}

	if theme != nil {
		viper.Set("theme", *theme)
		c.AuthCfg.Theme = *theme
	}

	if err := viper.WriteConfig(); err != nil {
		return err
	}

	return nil
}

func (c *cfg) HasEnvToken() bool {
	return c.AuthCfg.Token != "" || c.AuthCfg.Username != "" || c.AuthCfg.ProfileID != ""
}

func (c *cfg) DeleteConfig() error {
	c.AuthCfg.Username = ""
	c.AuthCfg.Role = ""
	c.AuthCfg.Token = ""
	c.AuthCfg.RefreshToken = ""
	c.AuthCfg.TeamId = ""
	c.AuthCfg.ProfileID = ""
	c.AuthCfg.EndPoint = "https://api.logfire.ai/"
	c.AuthCfg.GrpcEndpoint = "api.logfire.ai:443"
	c.AuthCfg.GrpcIngestion = "https://in.logfire.ai"

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
