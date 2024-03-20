package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Email    EmailConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Host string
	Port string
}

type EmailConfig struct {
	APIKey                         string
	FromEmail                      string
	SignupVerificationTemplateName string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
	PgDriver string
}

func LoadConfigFromFile(path string) (*viper.Viper, error) {
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

func ParseConfig(v *viper.Viper) (*Config, error) {
	var c Config

	err := v.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
