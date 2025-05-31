package strategy

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wer14/messenger/services/gateway/internal/app/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

var _ Strategy = (*HTTPServerStrategy)(nil)

type HTTPServerStrategy struct {
	httpPort              string
	httpServer            *http.Server
	httpGateway           *runtime.ServeMux
	httpGatewayRegistrars []server.HTTPGatewayRegistrar
	gatewayConnection     *grpc.ClientConn
}

func NewHTTPServerStrategy(
	httpPort string,
	httpServer *http.Server,
	httpGatewayRegistrars []server.HTTPGatewayRegistrar,
	httpGateway *runtime.ServeMux,
	gatewayConnection *grpc.ClientConn,
) *HTTPServerStrategy {
	return &HTTPServerStrategy{
		httpPort:              httpPort,
		httpServer:            httpServer,
		httpGateway:           httpGateway,
		httpGatewayRegistrars: httpGatewayRegistrars,
		gatewayConnection:     gatewayConnection,
	}
}

func (h *HTTPServerStrategy) Start() error {
	if h.gatewayConnection != nil { // may be nil if no grpc server needed
		if err := dialGRPCServer(h.gatewayConnection); err != nil {
			return fmt.Errorf("failed to dial grpc server: %w", err)
		}
	}

	if err := h.registerHTTPGateway(); err != nil {
		return fmt.Errorf("failed to register http gateway: %w", err)
	}

	h.startHTTPServer()

	return nil
}

func (h *HTTPServerStrategy) Stop() error {
	if err := h.httpServer.Close(); err != nil {
		return fmt.Errorf("failed to close http server: %w", err)
	}

	if h.gatewayConnection != nil {
		if err := h.gatewayConnection.Close(); err != nil {
			return fmt.Errorf("failed to close grpc client connection: %w", err)
		}
	}

	slog.Info("HTTP server stopped")

	return nil
}

func (h *HTTPServerStrategy) registerHTTPGateway() error {
	ctx := context.Background()
	// may be empty for http strategy
	for _, httpGatewayRegistrar := range h.httpGatewayRegistrars {
		if err := httpGatewayRegistrar.RegisterHTTPGateway(ctx, h.httpGateway, h.gatewayConnection); err != nil {
			return fmt.Errorf("failed to register http gateway: %w", err)
		}
	}

	h.httpServer.Handler = h.httpGateway

	return nil
}

func (h *HTTPServerStrategy) startHTTPServer() {
	go func() {
		if err := h.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("failed to listen and serve http server: %v", err)
		}
	}()

	slog.Info(fmt.Sprintf("serving grpc-gateway on http://localhost%s", h.httpPort))
	//slog.Info("serving swagger-ui on http://localhost%s/swagger-ui/", h.httpPort)
}

func dialGRPCServer(conn *grpc.ClientConn) error {
	waitCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	for {
		state := conn.GetState()

		if state == connectivity.Ready {
			break
		}

		if state == connectivity.Idle {
			conn.Connect()
		}

		if !conn.WaitForStateChange(waitCtx, state) {
			return fmt.Errorf("failed to establish connection: state did not change from %s", state.String())
		}

		if waitCtx.Err() != nil {
			return fmt.Errorf("failed to establish connection: context timeout %w", waitCtx.Err())
		}
	}

	return nil
}
