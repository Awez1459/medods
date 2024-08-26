package grpc_service_clients

import (
	"fmt"

	pbu "medods/api-service/genproto/user-proto"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"

	"medods/api-service/internal/pkg/config"
)

type ServiceClient interface {
	UserService() pbu.UserServiceClient
	Close()
}

type serviceClient struct {
	connections []*grpc.ClientConn
	userService pbu.UserServiceClient
}

func New(cfg *config.Config) (ServiceClient, error) {
	// user service

	connUserService, err := grpc.Dial(
		fmt.Sprintf("%s%s", cfg.UserService.Host, cfg.UserService.Port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	)
	if err != nil {
		return nil, err
	}

	return &serviceClient{
		userService: pbu.NewUserServiceClient(connUserService),
		connections: []*grpc.ClientConn{
			connUserService,
		},
	}, nil
}

func (s *serviceClient) UserService() pbu.UserServiceClient {
	return s.userService
}

func (s *serviceClient) Close() {
	for _, conn := range s.connections {
		if err := conn.Close(); err != nil {
			// should be replaced by logger soon
			fmt.Printf("error while closing grpc connection: %v", err)
		}
	}
}
