package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/wer14/messenger/services/auth/internal/server/grpc"
	"github.com/wer14/messenger/services/auth/internal/server/http"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Servers struct {
	grpcServer *grpc.GRPCServer
	httpServer *http.Server
}

func Start() (*Servers, error) {
	status := NewAtomicHealth()

	grpcServer, err := grpc.NewGRPCServer(":8082")
	if err != nil {
		return nil, fmt.Errorf("failed to init gRPC server: %w", err)
	}

	healthServer := health.NewServer()
	healthServer.SetServingStatus("auth.AuthService", healthpb.HealthCheckResponse_NOT_SERVING)
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	httpListener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, fmt.Errorf("failed to listen on HTTP port: %w", err)
	}

	httpServer := http.NewHTTPServer(":8080", status)
	go func() {
		if err := httpServer.Serve(httpListener); err != nil {
			log.Fatalf("HTTP debug server error: %v", err)
		}
	}()

	status.SetReady(true)
	healthServer.Resume()

	return &Servers{
		grpcServer: grpcServer,
		httpServer: httpServer,
	}, nil
}

func (s *Servers) Shutdown(ctx context.Context) error {
	log.Println("shutting down gRPC server...")
	s.grpcServer.Shutdown()

	log.Println("shutting down HTTP debug server...")
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Printf("error during HTTP shutdown: %v", err)
	}

	log.Println("shutdown complete.")
	return nil
}
