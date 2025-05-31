package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type GRPCRegistrar interface {
	RegisterServer(gs grpc.ServiceRegistrar)
}

type HTTPGatewayRegistrar interface {
	RegisterHTTPGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

type Server interface {
	Serve() error
	Stop() error
}
