package v2

import (
	"github.com/Reholly/kforge-proto/src/gen/auth"
	"google.golang.org/grpc"
	"net"
	"sso-service/api/http/v2/handler"
	"sso-service/internal/service"
)

type GrpcServer struct {
	server  *grpc.Server
	service service.ServiceManager
}

func NewGrpcServer(service service.ServiceManager) GrpcServer {
	return GrpcServer{
		service: service,
		server:  grpc.NewServer(),
	}
}

func (server *GrpcServer) Run(address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	authHandler := handler.NewAuthGrpcServer(server.service)
	auth.RegisterAuthServiceServer(server.server, authHandler)

	if err := server.server.Serve(listener); err != nil {
		return err
	}

	return nil
}

func (server *GrpcServer) Shutdown() {
	server.server.GracefulStop()
}
