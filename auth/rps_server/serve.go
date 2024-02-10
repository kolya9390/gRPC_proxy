package rpcserver

import (
	"auth/app"
	"auth/config"
	pb "auth/protoc/gen/auth"
	"auth/service/user"
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AurhService struct {
	AurhServiceGRPC
}

type AurhServiceGRPC struct {
	auther app.Auther
	pb.UnimplementedAuthServer
	
}

func NewAurhServis() *AurhService {
	return &AurhService{}
}

func (as *AurhService) StartServer() error {

	config := config.NewAppConf("server_app/.env")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", config.RPCServer.Port))
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()

	log.Printf("RPC типа %s сервер запущен и прослушивает порт :%s", config.RPCServer.Type, config.RPCServer.Port)
	//
	user_Service := user.NewGeoClient()

	as.auther = app.NewAuthProvider(user_Service)

	//
	grpcServer := grpc.NewServer()
	pb.RegisterAuthServer(grpcServer,
		&as.AurhServiceGRPC)
	grpcServer.Serve(listen)

	return nil
}

func (s *AurhServiceGRPC) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	token, err := s.auther.Login(ctx, in.GetEmail(), in.GetPassword())
	if err != nil {

		return nil, err
		//status.Error(codes.Internal, "failed to login")
	}

	return &pb.LoginResponse{Token: token}, err
}

func (s *AurhServiceGRPC) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	userID, err := s.auther.RegisterNewUser(ctx,in.Name,in.Email,in.Password)

	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{UserId: userID}, nil
}
