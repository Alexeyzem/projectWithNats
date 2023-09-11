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
	storage StorageConfig
}

func GetConfig() Config {
	return Config{}
}
