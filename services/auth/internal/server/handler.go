package server

import (
	"context"

	pb "github.com/wer14/messenger/services/auth/internal/gen/auth"
)

var _ pb.AuthServiceServer = (*AuthHandler)(nil)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) Register(ctx context.Context, request *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		UserId: "Register",
	}, nil
}
func (h *AuthHandler) Login(ctx context.Context, request *pb.LoginRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		UserId: "Login",
	}, nil
}
func (h *AuthHandler) OAuthLogin(ctx context.Context, request *pb.OAuthLoginRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		UserId: "OAuthLogin",
	}, nil
}
func (h *AuthHandler) VerifyEmail(ctx context.Context, request *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	return &pb.VerifyEmailResponse{
		Success: true,
	}, nil
}
