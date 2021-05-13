//+build wireinject

package main

import (
	_ "github.com/google/subcommands"
	"github.com/google/wire"

	"tiniapp-backend-oauth-sample/internal/repos"
	api_service "tiniapp-backend-oauth-sample/internal/services/api"

	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
)

func InitializeApplication() (*Application, func(), error) {
	wire.Build(
		ProvideConfig,
		ProvideLoggerConfig,
		pkg_logger.ProvideLogger,
		repos.NewRepository,
		ProvideAPIServiceConfig,
		api_service.NewService,
		ProvideApplication,
	)
	return &Application{}, nil, nil
}
