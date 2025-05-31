package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCRegistrar interface {
	RegisterServer(gs grpc.ServiceRegistrar)
}

type HTTPGatewayRegistrar interface {
	RegisterHTTPGateway(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error
}

type netServer interface {
	Serve() error
	Stop() error
}

type GRPCServer interface {
	netServer
	reflection.GRPCServer
}

type HTTPServer interface {
	netServer
}
