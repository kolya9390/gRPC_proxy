package main

import (
	"log"

	rpcserver "github.com/kolya9390/gRPC_GeoProvider/server_rpc/rps_server"
)

func main() {

	rpcServer := rpcserver.NewGeoServis()
	err := rpcServer.StartServer()

	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
