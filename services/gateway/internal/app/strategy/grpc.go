package strategy

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/wer14/messenger/services/gateway/internal/app/server"
	"google.golang.org/grpc/reflection"
)

var _ Strategy = (*GrpcStrategy)(nil)

type GrpcStrategy struct {
	grpcAddr       string
	grpcServer     server.GRPCServer
	grpcRegistrars []server.GRPCRegistrar
}

func NewGrpcServerStrategy(
	grpcAddr string,
	grpcServer server.GRPCServer,
	grpcRegistrars []server.GRPCRegistrar,
) *GrpcStrategy {
	return &GrpcStrategy{
		grpcAddr:       grpcAddr,
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
	if err := g.grpcServer.Stop(); err != nil {
		return fmt.Errorf("failed to stop gRPC server: %w", err)
	}

	slog.Info("gRPC server stopped")

	return nil
}

func (g *GrpcStrategy) startGRPCServer() error {
	go func() {
		if err := g.grpcServer.Serve(); err != nil {
			log.Fatalf("failed to serve grpc server: %v", err)
		}
	}()

	slog.Info(fmt.Sprintf("serving grpc-server on http://localhost%s", g.grpcAddr))

	return nil
}

func (g *GrpcStrategy) registerGRPCServices() {
	for _, grpcRegistrar := range g.grpcRegistrars {
		grpcRegistrar.RegisterServer(g.grpcServer)
	}

	reflection.Register(g.grpcServer)
}
