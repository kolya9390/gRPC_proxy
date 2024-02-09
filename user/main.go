package main

import (
	"log"

	rpcserver "user/rps_server"
)

func main() {

	rpcServer := rpcserver.NewUserServis()
	err := rpcServer.StartServer()

	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
