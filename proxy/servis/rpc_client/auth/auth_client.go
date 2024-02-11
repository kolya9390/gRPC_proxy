package auth_client

import (
	"context"
	"log"

	authService "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/protoc/gen/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcAuthClient struct {
	client *grpc.ClientConn
}

func NewUserClient() Auther {

	client, err := grpc.Dial("server_rpc:1237", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %s", err)
	}
	return &grpcAuthClient{client: client}

}

func (a *grpcAuthClient) Login(ctx context.Context, email, password string) (string, error){

	clientAuth := authService.NewAuthClient(a.client)

	resp,err := clientAuth.Login(ctx,&authService.LoginRequest{Email: email,Password: password})

	if err != nil {
		return "err login",err
	}

	return resp.GetToken(),nil

}

func (a *grpcAuthClient) Registeretion(ctx context.Context, name, email, password string) (int64, error){

	clientAuth := authService.NewAuthClient(a.client)

	resp,err := clientAuth.Register(ctx,&authService.RegisterRequest{Name: name,Email: email,Password: password})

	if err != nil {
		return 99,err
	}

	return resp.GetUserId(),nil

}
