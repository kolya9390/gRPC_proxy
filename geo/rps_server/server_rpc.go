package rpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/app"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/config"
	geo_provider "github.com/kolya9390/gRPC_GeoProvider/server_rpc/gen"
	servicegeo "github.com/kolya9390/gRPC_GeoProvider/server_rpc/servis_geo"
	"github.com/kolya9390/gRPC_GeoProvider/server_rpc/storage"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GeoServiceRPC struct {
	geoProvider app.GeoProvider
}

type GeoService struct {
	GeoServiceRPC
	GeoProviderGRPCServer
}

// GeoProviderGRPCServer must be embedded to have forward compatible implementations.
type GeoProviderGRPCServer struct {
	geoProvider app.GeoProvider
	geo_provider.UnimplementedGeoProviderGRPCServer
}

func RegisterGRPC(gRPC *grpc.Server) {
	geo_provider.RegisterGeoProviderGRPCServer(gRPC, &GeoProviderGRPCServer{})
}

func NewGeoServis() *GeoService {
	return &GeoService{
	}
}

func (gs *GeoService) StartServer(port string) error {

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

	postgresDB := storage.NewGeoRepositoryDB(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Cache.Address, config.Cache.Port),
	})

	defer db.Close()
	defer redisClient.Close()

	cache := storage.NewGeoRedis(redisClient)
	storageDB := storage.NewGeoRepositoryProxy(*postgresDB, cache)
	sevisDAdata := servicegeo.NewDadataService(config.AuthorizationDADATA)

	err = postgresDB.ConnectToDB()

	if err != nil {
		log.Printf("Error conect DB %s", err)
	}

	gs.GeoServiceRPC.geoProvider = app.NewGeoProvider(storageDB,sevisDAdata)
	gs.GeoProviderGRPCServer.geoProvider = app.NewGeoProvider(storageDB,sevisDAdata)
	
	//

	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", config.RPCServer.Port))
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()



	log.Printf("RPC типа %s сервер запущен и прослушивает порт :%s", config.RPCServer.Type, config.RPCServer.Port)

	//
		grpcServer := grpc.NewServer()
		geo_provider.RegisterGeoProviderGRPCServer(grpcServer, 
			&gs.GeoProviderGRPCServer)
		grpcServer.Serve(listen)
		

	return nil
}

func (gs *GeoServiceRPC) AddressSearchRPC(query RequestAddressSearch, reply *[]*Address) error {
	addresses, err := gs.geoProvider.AddressSearch(query.Query)
	if err != nil {
		log.Printf("Error AddressSearch: %v", err)
		return err
	}

	for _, adres := range addresses {
		*reply = append(*reply, &Address{
			GeoLat: adres.GeoLat,
			GeoLon: adres.GeoLon,
			Result: adres.Result,
		})

	}

	return nil
}

func (gs *GeoServiceRPC) AddressGeoCodeRPC(geocode RequestAddressGeocode, reply *[]*Address) error {
	addresses, err := gs.geoProvider.GeoCode(geocode.Lat, geocode.Lng)
	if err != nil {
		log.Printf("Error AddressGeoCode: %v", err)
		return err
	}
	// Просто присваиваем новое значение reply через косвенное разыменование
	for _, adres := range addresses {
		*reply = append(*reply, &Address{
			GeoLat: adres.GeoLat,
			GeoLon: adres.GeoLon,
			Result: adres.Result,
		})

	}

	return nil
}

func (gs *GeoProviderGRPCServer) AddressSearchGRPC(ctx context.Context, req *geo_provider.RequestAddressSearch) (*geo_provider.RespAddress, error) {

	addresses, err := gs.geoProvider.AddressSearch(req.Query)
	if err != nil {
		log.Printf("Error AddressSearch: %v", err)
		return nil, err
	}

	return &geo_provider.RespAddress{
		Geolat: addresses[0].GeoLat,
		Geolon: addresses[0].GeoLon,
		Result: addresses[0].Result,
	}, nil
}
func (gs *GeoProviderGRPCServer) AddressGeoCodeGRPC(ctx context.Context, req *geo_provider.RequestAddressGeocode) (*geo_provider.RespAddress, error) {

	addresses, err := gs.geoProvider.GeoCode(req.Lat, req.Lng)
	if err != nil {
		log.Printf("Error AddressGeoCode: %v", err)
		return nil, err
	}

	return &geo_provider.RespAddress{
		Geolat: addresses[0].GeoLat,
		Geolon: addresses[0].GeoLon,
		Result: addresses[0].Result,
	}, status.Errorf(codes.Unimplemented, "method AddressGeoCodeGRPC not implemented")
}
