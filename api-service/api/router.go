package api

import (
	// "net/http"
	"time"

	_ "medods/api-service/api/docs"
	v1 "medods/api-service/api/handlers/v1"

	"medods/api-service/api/middleware"

	"github.com/casbin/casbin/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"

	grpcClients "medods/api-service/internal/infrastructure/grpc_service_client"
	"medods/api-service/internal/pkg/config"
	tokens "medods/api-service/internal/pkg/token"
	"medods/api-service/internal/usecase/app_version"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	JwtHandler     tokens.JwtHandler
	AppVersion     app_version.AppVersion
	Enforcer       *casbin.Enforcer
}

// NewRouter
// @title Welcome To Booking API
// @Description API for Touristan
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Service:        option.Service,
		JwtHandler:     option.JwtHandler,
		AppVersion:     option.AppVersion,
		Enforcer:       option.Enforcer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	// router.Use(middleware.Tracing)
	router.Use(middleware.CheckCasbinPermission(option.Enforcer, *option.Config))
	router.Static("/media", "./media")
	api := router.Group("/v1")

	// AUTH METHODS
	api.POST("/users/login", HandlerV1.Token)
	api.GET("/token/:refresh", HandlerV1.UpdateToken)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	return router
}
