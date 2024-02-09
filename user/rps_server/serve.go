package rpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
	"user/app"
	"user/config"
	userApp "user/protoc/gen"
	"user/storage"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserService struct {
	UserServiceGRPC
}

type UserServiceGRPC struct {
	userProvider app.UserProvider
	userApp.UnimplementedUserServiceServer
}

func NewUserServis() *UserService {
	return &UserService{}
}

func (us *UserService) StartServer() error {

	config := config.NewAppConf("server_app/.env")

	// Инициализация подключения к базе данных
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

	db, err := sqlx.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}
	time.Sleep(time.Second * 3)
	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %s", err)
	}

	defer db.Close()

	postgresDB := storage.NewUserRepositoryDB(db)

	us.userProvider = app.NewUserProvider(postgresDB)

	err = postgresDB.ConnectToDB()

	if err != nil {
		log.Printf("Error conect DB %s", err)
	}

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", config.RPCServer.Port))
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()

	log.Printf("RPC типа %s сервер запущен и прослушивает порт :%s", config.RPCServer.Type, config.RPCServer.Port)

	//
	grpcServer := grpc.NewServer()
	userApp.RegisterUserServiceServer(grpcServer,
		&us.UserServiceGRPC)
	grpcServer.Serve(listen)

	return nil
}

func (us *UserServiceGRPC) GetUserProfileIDs(ctx context.Context, user_id *userApp.RequestUserID) (*userApp.ResponseUser, error) {

	user, err := us.userProvider.GetUserIDs(user_id.GetUserId())

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &userApp.ResponseUser{
		UserId: user.ID, 
		Name: user.Name, 
		Email: user.Email,
		Password: string(user.Password)}, nil

}

func (us *UserServiceGRPC) GetListUser(ctx context.Context,empty *userApp.Empty) (*userApp.ResponseListUsers, error) {

	respUsers := &userApp.ResponseListUsers{}

	users, err := us.userProvider.GetAllUser()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	for _, user := range users {
		respUsers.Users = append(respUsers.Users,
			&userApp.ResponseUser{
				UserId: user.ID,
				Name:   user.Name,
				Email:  user.Email,
				Password: string(user.Password),
			})
	}

	return respUsers, nil

}

func (us *UserServiceGRPC) Register(ctx context.Context, in *userApp.RegisterRequest) (*userApp.RegisterResponse, error) {

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	userID, err := us.userProvider.RegisterNewUser(ctx, in.GetName(), in.GetEmail(), in.GetPassword())

	if err != nil {
		return nil, err
	}

	return &userApp.RegisterResponse{UserId: userID}, nil
}
