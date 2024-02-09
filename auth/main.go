package main

import (
//	"fmt"
	"log"
//	"time"
	_ "github.com/lib/pq"
	"auth/app"
	//"auth/config"
	rpcserver "auth/rps_server"
	"auth/storage"

	"github.com/jmoiron/sqlx"
)

func main() {

	//config := config.NewAppConf("server_app/.env")
/*
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
*/
	postgresDB := storage.NewAuthRepositoryDB(&sqlx.DB{})

	appAuthService := app.NewAuthApp(*postgresDB,10)

	rpcServer := rpcserver.NewAuthService(appAuthService,1237)
	err := rpcServer.StartServer()

	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
