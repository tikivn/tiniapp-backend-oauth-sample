package main

import (
	"sync"

	_ "go.uber.org/automaxprocs"

	api_service "tiniapp-backend-oauth-sample/internal/services/api"

	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
	pkg_signal "tiniapp-backend-oauth-sample/pkg/signal"
)

type Application struct {
	Config     *Config
	Logger     pkg_logger.ILogger
	APIService *api_service.Service
}

func (a *Application) Start() {
	logger := pkg_logger.GetLogger().WithPrefix("Application")
	logger.Info("app is starting")

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		go a.APIService.Start()
	}(&wg)

	wg.Wait()
}

func (a *Application) Stop() {
	logger := pkg_logger.GetLogger().WithPrefix("Application")

	logger.Info("app is stopping")
	// TODO ...
	logger.Info("app is stopped!")
}

func NewApplication(
	config *Config,
	logger pkg_logger.ILogger,
	restAPIService *api_service.Service,
) (*Application, error) {
	app := &Application{
		Config:     config,
		Logger:     logger,
		APIService: restAPIService,
	}
	return app, nil
}

func ProvideApplication(
	config *Config,
	logger pkg_logger.ILogger,
	restAPIService *api_service.Service,
) (*Application, func(), error) {
	app, err := NewApplication(
		config,
		logger,
		restAPIService,
	)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		app.Logger.Info("app is cleaning")
		app.Stop()
		app.Logger.Info("app is cleaned!")
	}
	return app, cleanup, nil
}

func main() {
	logger := pkg_logger.GetLogger().WithPrefix("Application")

	app, cleanup, err := InitializeApplication()
	if err != nil {
		app.Logger.Panic(err)
		panic(err)
	}
	defer cleanup()

	app.Start()

	signal := pkg_signal.WaitOSSignal()
	logger.Infof("received signal: %s", signal.String())
}
