package strategy

import (
	"fmt"
	"log"
	"log/slog"
	"net"

	"github.com/wer14/messenger/services/gateway/internal/app/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var _ Strategy = (*GrpcStrategy)(nil)

type GrpcStrategy struct {
	grpcPort       string
	grpcServer     *grpc.Server
	grpcRegistrars []server.GRPCRegistrar
}

func NewGrpcServerStrategy(
	grpcPort string,
	grpcServer *grpc.Server,
	grpcRegistrars []server.GRPCRegistrar,
) *GrpcStrategy {
	return &GrpcStrategy{
		grpcPort:       grpcPort,
		grpcServer:     grpcServer,
		grpcRegistrars: grpcRegistrars,
	}
}

func (g *GrpcStrategy) Start() error {
	g.registerGRPCServices()

	if err := g.startGRPCServer(); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}

	return nil
}

func (g *GrpcStrategy) Stop() error {
	g.grpcServer.GracefulStop()
	slog.Info("gRPC server stopped")

	return nil
}

func (g *GrpcStrategy) startGRPCServer() error {
	listener, err := net.Listen("tcp", g.grpcPort)
	if err != nil {
		return fmt.Errorf("failed to listen on port %s: %w", g.grpcPort, err)
	}

	go func() {
		if err := g.grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve grpc server: %v", err)
		}
	}()

	slog.Info("serving grpc-server on http://localhost%s", g.grpcPort)

	return nil
}

func (g *GrpcStrategy) registerGRPCServices() {
	for _, grpcRegistrar := range g.grpcRegistrars {
		grpcRegistrar.RegisterServer(g.grpcServer)
	}

	reflection.Register(g.grpcServer)
}
