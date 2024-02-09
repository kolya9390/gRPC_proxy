package rpcserver

import (
	"auth/app"
	grpcauth "auth/rps_server/grpcAuth"
	"context"
	"log"
	"net"
	"net/rpc"

	"google.golang.org/grpc"
)



func NewAuthService(authservice app.Auther, port int) *AuthService {

	gRPCServer := grpc.NewServer()

	grpcauth.Register(gRPCServer, authservice)

	return &AuthService{gRPCServer: gRPCServer}
}

func (as *AuthService) Stop() {
	const op = "grpcapp.Stop"

	log.Println("stop")

	as.gRPCServer.GracefulStop()
}

func (as *AuthService) StartServer() error {


	listen, err := net.Listen("tcp",":1237")
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()

	log.Printf("%s", listen)

	log.Println("RPC сервер запущен и прослушивает порт :1237")
	rpc.Accept(listen)

	return nil
}

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
