package main

import (
	"log"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wer14/messenger/services/gateway/internal/app"
	"github.com/wer14/messenger/services/gateway/internal/app/server"
	"github.com/wer14/messenger/services/gateway/internal/app/strategy"
	"github.com/wer14/messenger/services/gateway/internal/handler"
	"github.com/wer14/messenger/services/gateway/internal/services/gateway"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpAddr string = ":8080"
	grpcPort string = "50051"
)

func main() {
	gatewayServer := initGatewayServer()
	gatewayClient := initGRPCClient(grpcPort)

	grpcGatewayMux := runtime.NewServeMux()
	handler.RegisterHealthSystemRoutes(grpcGatewayMux)

	app := app.NewApp(
		strategy.NewGrpcServerStrategy(
			grpcPort,
			newGRPCServer(),
			[]server.GRPCRegistrar{gatewayServer},
		),
		strategy.NewHTTPServerStrategy(
			httpAddr,
			newHTTPServer(httpAddr),
			[]server.HTTPGatewayRegistrar{gatewayServer},
			grpcGatewayMux,
			gatewayClient,
		),
	)

	if err := app.Run(); err != nil {
		log.Fatalf("infra start error: %v", err)
	}

	app.Stop()
}

func initGatewayServer() *gateway.Server {
	gatewayServer := gateway.NewGatewayServer()

	return gatewayServer
}

func newHTTPServer(httpPort string) *http.Server {
	httpServer := &http.Server{
		Addr:         httpPort,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}

func newGRPCServer() *grpc.Server {
	grpcServer := grpc.NewServer()

	return grpcServer
}

func initGRPCClient(grpcPort string) *grpc.ClientConn {
	conn, err := grpc.NewClient(
		"127.0.0.1"+grpcPort,
		[]grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}...,
	)
	if err != nil {
		log.Fatalf("failed to create grpc client connection: %v", err)
	}

	return conn
}
