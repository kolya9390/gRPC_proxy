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
)

type UserService struct {
	UserServiceGRPC
}

type UserServiceGRPC struct {
	userProvider app.UserProvider
	userApp.UnimplementedUserServiceServer
}

func NewGeoServis() *UserService{
	return &UserService{}
}

func (us *UserService) StartServer(port string) error {

	config := config.NewAppConf("server_app/.env")

	// Инициализация подключения к базе данных
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)

		log.Fatal(connstr)

	db, err := sqlx.Open("postgres", connstr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %s", err)
	}

	time.Sleep(time.Second * 3)
	// Проверка соединения с базой данных
	if err := db.Ping(); err != nil {
		log.Printf("Error pinging the database: %s", err)
	}

	defer db.Close()

	postgresDB := storage.NewGeoRepositoryDB(db)

	us.userProvider = app.NewGeoProvider(postgresDB)

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


	grpcServer := grpc.NewServer()
		userApp.RegisterUserServiceServer(grpcServer, 
			&us.UserServiceGRPC)
		grpcServer.Serve(listen)

	return nil
}

func (us *UserService) GetUserProfileIDs(ctx context.Context,user_id *userApp.RequestUserID) (*userApp.ResponseUser, error)  {

	user , err := us.userProvider.GetUserIDs(user_id.UserId)

	if err != nil {
		log.Println(err)
		return nil,err
	}

	return &userApp.ResponseUser{Email: user.Email,Password: user.Password},nil
    
}


func (us *UserService) GetListUser(context.Context, *userApp.Empty) (*userApp.ResponseListUsers, error)  {

	var respUsers *userApp.ResponseListUsers

	users, err := us.userProvider.GetAllUser()
	if err != nil {
		log.Println(err)
		return nil,err
	}

	for _, user := range users{
		respUsers.Users = append(respUsers.Users, 
			&userApp.ResponseUser{
				Email: user.Email,
				Password: user.Password,
				}) 
	}

	return respUsers,nil
		
	}
