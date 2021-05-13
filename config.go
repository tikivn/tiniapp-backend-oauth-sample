package main

import (
	"log"

	api_service "tiniapp-backend-oauth-sample/internal/services/api"

	pkg_cfgutil "tiniapp-backend-oauth-sample/pkg/cfgutil"
	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
)

type Config struct {
	Environment string              `yaml:"environment"`
	Logger      *pkg_logger.Config  `yaml:"logger"`
	APIService  *api_service.Config `yaml:"api_service"`
}

func LoadConfig() *Config {
	var config Config
	err := pkg_cfgutil.LoadConfig("", &config)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return &config
}

func ProvideConfig() *Config {
	return LoadConfig()
}

func ProvideLoggerConfig(config *Config) *pkg_logger.Config {
	return config.Logger
}

func ProvideAPIServiceConfig(config *Config) *api_service.Config {
	return config.APIService
}
