package gateway

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/wer14/messenger/services/gateway/internal/app/server"
	pb "github.com/wer14/messenger/services/gateway/internal/gen/gateway"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

var (
	_ server.GRPCRegistrar        = (*Server)(nil)
	_ server.HTTPGatewayRegistrar = (*Server)(nil)
)

type Server struct {
	ValidateMessages []proto.Message
	ValidatePaths    map[string]bool

	pb.UnimplementedGatewayServiceServer
}

func NewGatewayServer() *Server {
	return &Server{
		ValidateMessages: []proto.Message{
			&pb.HelloRequest{},
		},
		ValidatePaths: map[string]bool{
			"/gateway.GatewayService/Hello": true,
		},
	}
}

func (s *Server) RegisterServer(gs grpc.ServiceRegistrar) {
	pb.RegisterGatewayServiceServer(gs, s)
}

func (s *Server) RegisterHTTPGateway(
	ctx context.Context,
	mux *runtime.ServeMux,
	conn *grpc.ClientConn,
) error {
	if err := pb.RegisterGatewayServiceHandler(ctx, mux, conn); err != nil {
		return fmt.Errorf("failed to register gateway api handler: %w", err)
	}

	return nil
}

func (s *Server) Hello(context.Context, *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{
		Msg: "HELLO!",
	}, nil
}
