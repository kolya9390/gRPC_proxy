package user

import (
	"auth/config"
	userApp "auth/protoc/gen/user"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type grpcUserClient struct {
	client *grpc.ClientConn
}

func NewGeoClient() UserService {

	config := config.NewAppConf("server_app/.env")

	log.Println(config.RPCServer.Port)

	client, err := grpc.Dial("server_rpc:1235", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect %s", err)
	}
	return &grpcUserClient{client: client}

}

func (u *grpcUserClient) GetUserIDs(ctx context.Context, user_id int64) (User, error) {
	clientUser := userApp.NewUserServiceClient(u.client)

	resp, err := clientUser.GetUserProfileIDs(ctx, &userApp.RequestUserID{UserId: user_id})

	if err != nil {
		return User{},err

	}

	return User{
		ID: resp.UserId,
		Name: resp.Name,
		Email: resp.Email,
		Password: []byte(resp.Password),},nil
}


func (u *grpcUserClient) GetListUsers(ctx context.Context) ([]User,error){

	clientUser := userApp.NewUserServiceClient(u.client)

	var users []User

	resp , err := clientUser.GetListUser(ctx,&userApp.Empty{})

	
	if err != nil {
		return nil,err

	}

	for _,user := range resp.Users {
		users = append(users, User{
			ID: user.UserId,
			Name: user.Name,
			Email: user.Email,
			Password: []byte(user.Password),
		})
	}

	return users,nil

}

func (u *grpcUserClient) Register(ctx context.Context,name, email, password string) (int64, error) {

	clientUser := userApp.NewUserServiceClient(u.client)

	resp,err := clientUser.Register(ctx,&userApp.RegisterRequest{
		Name: name,
		Email: email,
		Password: password,
	})

	if err != nil {
		return 0, err
	}

	return resp.UserId,err

}
