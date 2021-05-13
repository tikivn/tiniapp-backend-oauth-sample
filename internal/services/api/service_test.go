package api_test

import (
	"os"
	"testing"

	pkg_cfgutil "tiniapp-backend-oauth-sample/pkg/cfgutil"
	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"

	"tiniapp-backend-oauth-sample/internal/repos"
	api_service "tiniapp-backend-oauth-sample/internal/services/api"
)

type testConfig struct {
	Environment string              `yaml:"environment"`
	Logger      *pkg_logger.Config  `yaml:"logger"`
	APIService  *api_service.Config `yaml:"api_service"`
}

var svc *api_service.Service

func TestMain(m *testing.M) {
	logger := pkg_logger.GetLogger()

	exitVal := int(1)

	var cfg testConfig
	err := pkg_cfgutil.LoadConfig("", &cfg)
	if err != nil {
		logger.Error(err)
		os.Exit(exitVal)
	}

	testRepo := repos.NewRepository()

	svc, err = api_service.NewService(cfg.APIService, testRepo)
	if err != nil {
		logger.Error(err)
		os.Exit(exitVal)
	}

	svc.ConfigureRoute()

	exitVal = m.Run()

	os.Exit(exitVal)
}
