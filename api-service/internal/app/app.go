package app

import (
	"context"
	"fmt"
	"medods/api-service/api"
	grpcService "medods/api-service/internal/infrastructure/grpc_service_client"
	"medods/api-service/internal/infrastructure/repository/postgresql"
	"medods/api-service/internal/pkg/config"
	"medods/api-service/internal/pkg/logger"
	"medods/api-service/internal/pkg/postgres"
	"medods/api-service/internal/pkg/redis"
	"medods/api-service/internal/usecase/app_version"
	"net/http"
	"time"

	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	defaultrolemanager "github.com/casbin/casbin/v2/rbac/default-role-manager"
	"go.uber.org/zap"
)

type App struct {
	Config     *config.Config
	Logger     *zap.Logger
	DB         *postgres.PostgresDB
	RedisDB    *redis.RedisDB
	server     *http.Server
	Enforcer   *casbin.Enforcer
	Clients    grpcService.ServiceClient
	appVersion app_version.AppVersion
}

func NewApp(cfg config.Config) (*App, error) {
	// logger init
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	// postgres init
	db, err := postgres.New(&cfg)
	if err != nil {
		return nil, err
	}

	// redis init
	redisdb, err := redis.New(&cfg)
	if err != nil {
		return nil, err
	}

	// initialization enforcer
	enforcer, err := casbin.NewEnforcer("auth.conf", "auth.csv")
	if err != nil {
		return nil, err
	}

	var contextTimeout time.Duration

	// context timeout initialization
	contextTimeout, err = time.ParseDuration(cfg.Context.Timeout)
	if err != nil {
		return nil, err
	}

	appVersionRepo := postgresql.NewAppVersionRepo(db)

	appVersionUseCase := app_version.NewAppVersionService(contextTimeout, appVersionRepo)

	return &App{
		Config:     &cfg,
		Logger:     logger,
		DB:         db,
		RedisDB:    redisdb,
		Enforcer:   enforcer,
		appVersion: appVersionUseCase,
	}, nil
}

func (a *App) Run() error {
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error while parsing context timeout: %v", err)
	}

	clients, err := grpcService.New(a.Config)
	if err != nil {
		return err
	}
	a.Clients = clients

	// api init
	handler := api.NewRoute(api.RouteOption{
		Config:         a.Config,
		Logger:         a.Logger,
		ContextTimeout: contextTimeout,
		Enforcer:       a.Enforcer,
		Service:        clients,
		AppVersion:     a.appVersion,
	})
	err = a.Enforcer.LoadPolicy()
	if err != nil {
		return err
	}
	roleManager := a.Enforcer.GetRoleManager().(*defaultrolemanager.RoleManagerImpl)

	roleManager.AddMatchingFunc("keyMatch", util.KeyMatch)
	roleManager.AddMatchingFunc("keyMatch3", util.KeyMatch3)

	// server init
	a.server, err = api.NewServer(a.Config, handler)
	if err != nil {
		return fmt.Errorf("error while initializing server: %v", err)
	}

	return a.server.ListenAndServe()
}

func (a *App) Stop() {

	// close database
	a.DB.Close()

	// close grpc connections
	a.Clients.Close()

	// shutdown server http
	if err := a.server.Shutdown(context.Background()); err != nil {
		a.Logger.Error("shutdown server http ", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
