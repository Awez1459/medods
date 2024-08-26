package v1

import (
	"time"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"

	grpcClients "medods/api-service/internal/infrastructure/grpc_service_client"
	"medods/api-service/internal/pkg/config"
	tokens "medods/api-service/internal/pkg/token"

	appV "medods/api-service/internal/usecase/app_version"
)

type HandlerV1 struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Service        grpcClients.ServiceClient
	AppVersion     appV.AppVersion
	Enforcer       *casbin.Enforcer
}

type HandlerV1Config struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	JwtHandler     tokens.JwtHandler
	Service        grpcClients.ServiceClient
	AppVersion     appV.AppVersion
	Enforcer       *casbin.Enforcer
}

func New(c *HandlerV1Config) *HandlerV1 {
	return &HandlerV1{
		Config:         c.Config,
		Logger:         c.Logger,
		ContextTimeout: c.ContextTimeout,
		Service:        c.Service,
		JwtHandler:     c.JwtHandler,
		AppVersion:     c.AppVersion,
		Enforcer:       c.Enforcer,
	}
}
