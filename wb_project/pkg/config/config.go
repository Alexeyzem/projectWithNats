package config

import "wb_project/pkg/logging"

type StorageConfig struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
	Lg       *logging.Logger
}

type Config struct {
	Storage StorageConfig
}

func GetConfig() *Config {
	return &Config{
		Storage: StorageConfig{
			Username: "alexey",
			Password: "alexey",
			Host:     "0.0.0.0",
			Port:     "8090",
			Database: "postgres",
		},
	}
}
