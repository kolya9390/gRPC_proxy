package grpcauth

import (
	authService "auth/protoc/gen"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	auth Auth
	authService.UnimplementedAuthServer
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	authService.RegisterAuthServer(gRPCServer, &ServerAPI{auth: auth})
}


func (s *ServerAPI) Login(ctx context.Context, in *authService.LoginRequest) (*authService.LoginResponse, error) {

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {


		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &authService.LoginResponse{Token: token},err
}


func (s *ServerAPI) RegisterNewUser(ctx context.Context, in *authService.RegisterRequest) (*authService.RegisterResponse, error) {

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, in.GetEmail(), in.GetPassword())

	if err != nil {
		return nil, err
	}

	return &authService.RegisterResponse{UserId: userID},nil
}
