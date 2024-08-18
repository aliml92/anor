package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/viper"
)

const (
	// ConfigEnvVar is the environment variable name for the config file path
	ConfigEnvVar = "CONFIG"

	// DefaultConfigPath is the default path for the config file
	DefaultConfigPath = "./config.yaml"
)

// Config holds all configuration for the application
type Config struct {
	Server      ServerConfig      `mapstructure:"server"`
	Email       EmailConfig       `mapstructure:"email"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Redis       RedisConfig       `mapstructure:"redis"`
	Session     SessionConfig     `mapstructure:"session"`
	Typesense   TypesenseConfig   `mapstructure:"typesense"`
	Stripe      StripeConfig      `mapstructure:"stripe"`
	GoogleOAuth GoogleOAuthConfig `mapstructure:"google_oauth"`
	Logger      LoggerConfig      `mapstructure:"logger"`
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Host    string
	Port    string
	IsHTTPS bool
}

// EmailConfig holds email-specific configuration
type EmailConfig struct {
	APIKey                         string
	FromEmail                      string
	Templates                      string
	SignupVerificationTemplateName string
}

// DatabaseConfig holds database-specific configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
	PgDriver string
}

// RedisConfig holds Redis-specific configuration
type RedisConfig struct {
	Host         string
	Port         string
	Username     string
	Password     string
	DB           int
	MinIdleConns int
	MaxIdleConns int
}

type SessionConfig struct {
	CookieName    string
	AuthLifetime  time.Duration
	GuestLifetime time.Duration
}

type TypesenseConfig struct {
	Host     string
	Port     int
	APIKey   string
	UseHTTPS bool
}

// StripeConfig holds Stripe-specific configuration
type StripeConfig struct {
	SecretKey string
}

type LoggerConfig struct {
	Level      string
	Format     string
	AddSource  bool
	TimeFormat string
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

// New creates and returns a new Config instance
func New() (*Config, error) {
	cfgPath := os.Getenv(ConfigEnvVar)
	if cfgPath == "" {
		cfgPath = DefaultConfigPath
	}

	cfgFile, err := loadConfigFromFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file: %v", err)
	}

	cfg, err := parseConfig(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return cfg, nil
}

func loadConfigFromFile(path string) (*viper.Viper, error) {
	v := viper.New()

	dir, file := filepath.Split(path)

	v.AddConfigPath(dir)
	v.SetConfigName(strings.TrimSuffix(file, filepath.Ext(file)))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, err
		}
		return nil, err
	}

	return v, nil
}

func parseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
