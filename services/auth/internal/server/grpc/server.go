package grpc

import (
	"fmt"
	"log"
	"net"

	"buf.build/go/protovalidate"
	pb "github.com/wer14/messenger/services/auth/internal/gen/auth"
	"github.com/wer14/messenger/services/auth/internal/server/grpc/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	*grpc.Server
	listener net.Listener
}

func NewGRPCServer(addr string) (*GRPCServer, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(),
		protovalidate.WithMessages(
			&pb.OAuthLoginRequest{},
			&pb.RegisterRequest{},
			&pb.LoginRequest{},
			&pb.VerifyEmailRequest{},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(
			interceptor.MethodScopedValidatioUnaryInterceptor(
				validator,
				map[string]bool{
					"/auth.AuthService/Register":    true,
					"/auth.AuthService/Login":       true,
					"/auth.AuthService/OAuthLogin":  true,
					"/auth.AuthService/VerifyEmail": true,
				},
			),
		),
	)
	reflection.Register(gRPCServer)

	pb.RegisterAuthServiceServer(gRPCServer, NewAuthHandler())

	return &GRPCServer{
		Server:   gRPCServer,
		listener: listener,
	}, nil
}

func (s *GRPCServer) Start() error {
	log.Printf("grpc server listening on %s", s.listener.Addr().String())

	return s.Server.Serve(s.listener)
}

func (s *GRPCServer) Shutdown() {
	log.Println("shutting down grpc server...")

	s.Server.GracefulStop()
}
