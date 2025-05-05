package server

import (
	"fmt"
	"log"
	"net"

	"buf.build/go/protovalidate"
	pb "github.com/wer14/messenger/services/auth/internal/gen/auth"
	"github.com/wer14/messenger/services/auth/internal/server/interceptor"
	"google.golang.org/grpc"
)

func Start(addr string) error {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
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
		return fmt.Errorf("failed to initialize validator: %w", err)
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.ValidationUnaryInterceptor(validator)),
	)

	h := NewAuthHandler()
	pb.RegisterAuthServiceServer(s, h)

	log.Printf("auth grpc server listening on %s", addr)
	return s.Serve(lis)
}
