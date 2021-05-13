package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	resty "github.com/go-resty/resty/v2"
	"moul.io/http2curl"

	"tiniapp-backend-oauth-sample/internal/repos"
	"tiniapp-backend-oauth-sample/internal/services/api/middlewares"

	pkg_logger "tiniapp-backend-oauth-sample/pkg/logger"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`

	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`

	TiniAppServerAddress string `yaml:"tiniapp_server_address"`
}

type Service struct {
	Config *Config

	Repository repos.IRepository

	httpRouter *gin.Engine
	httpServer *http.Server

	HTTPClient *resty.Client
}

func (s *Service) GetRouter() *gin.Engine {
	return s.httpRouter
}

func (s *Service) ConfigureRoute() {
	s.httpRouter.Use(middlewares.SetRequestID())

	s.httpRouter.GET("/api/ping", s.Ping())
	s.httpRouter.POST("/api/auth/me", s.GetMeFromAccessToken())
	s.httpRouter.POST("/api/auth/token", s.GetAccessTokenFromAuthCode())
	s.httpRouter.POST("/api/auth/token/refresh", s.GetAccessTokenFromRefreshToken())
}

func (s *Service) Start() {
	logger := pkg_logger.GetLogger().WithPrefix("Service")

	s.ConfigureRoute()

	logger.Infof("starting HTTP server at http://%v ...", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func (s *Service) Stop() {
	logger := pkg_logger.GetLogger().WithPrefix("Service")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := s.httpServer.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		logger.Fatal(err)
	}
}

func NewService(
	config *Config,
	repository repos.IRepository,
) (*Service, error) {
	httpRouter := gin.New()

	httpAddr := ":8080"
	if config.ServerAddress != "" {
		httpAddr = config.ServerAddress
	}

	httpServer := &http.Server{
		Addr:    httpAddr,
		Handler: httpRouter,
	}

	service := &Service{
		Config:     config,
		Repository: repository,
		httpRouter: httpRouter,
		httpServer: httpServer,
	}

	httpClient := resty.New().
		SetDebug(false).
		SetPreRequestHook(
			func(c *resty.Client, req *http.Request) error {
				logger := pkg_logger.GetLogger().WithPrefix("Service")

				command, _ := http2curl.GetCurlCommand(req)
				commandString := command.String()

				logger.Debug(commandString)

				return nil
			},
		)
	service.HTTPClient = httpClient

	return service, nil
}

func ProvideService(
	config *Config,
	repository repos.IRepository,
) (*Service, func(), error) {
	service, err := NewService(config, repository)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() {
		service.Stop()
	}
	return service, cleanup, nil
}
