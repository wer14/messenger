package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	grpcAddr string = ":8082"
)

func main() {
	gatewayServer := initGatewayServer()
	gatewayClient := initGRPCClient(grpcAddr)

	grpcGatewayMux := runtime.NewServeMux()
	handler.RegisterHealthSystemRoutes(grpcGatewayMux)

	app := app.NewApp(
		strategy.NewGrpcServerStrategy(
			grpcAddr,
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
		log.Fatalf("app run error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	slog.Info("stop signal received")
	if err := app.Stop(); err != nil {
		log.Fatalf("app stop failed: %v", err)
	}
}

func initGatewayServer() *gateway.Server {
	gatewayServer := gateway.NewGatewayServer()

	return gatewayServer
}

func newHTTPServer(httpAddr string) *http.Server {
	httpServer := &http.Server{
		Addr:         httpAddr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return httpServer
}

func newGRPCServer() server.GRPCServer {
	grpcServer, err := server.NewGRPCServer(server.GRPCParams{
		Address: grpcAddr,
	})
	if err != nil {
		log.Fatalf("failed to create grpc server: %v", err)
	}

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
