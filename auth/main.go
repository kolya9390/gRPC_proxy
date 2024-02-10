package main

import (
	"log"

	rpcserver "auth/rps_server"
)

func main() {

	rpcServer := rpcserver.NewAurhServis()
	err := rpcServer.StartServer()

	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}