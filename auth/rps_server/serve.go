package rpcserver

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"time"

	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/kolya9390/RPCGeoProvider/server_rpc/app"
	"github.com/kolya9390/RPCGeoProvider/server_rpc/config"
	servisgeo "github.com/kolya9390/RPCGeoProvider/server_rpc/servis_geo"
	"github.com/kolya9390/RPCGeoProvider/server_rpc/storage"
	_ "github.com/lib/pq"
)

type GeoService struct {
	geoProvider app.GeoProvider
}

func NewGeoServis() *GeoService{
	return &GeoService{}
}

func (gs *GeoService) StartServer(port string) error {

	env, err := godotenv.Read("server_app/.env")

	if err != nil {
		log.Println(err)
	}

	//	log.Println(env)

	config := &config.AppConf{
		DB: config.DB{
			Host:     env["DB_HOST"],
			Port:     env["DB_PORT"],
			User:     env["DB_USER"],
			Password: env["DB_PASSWORD"],
			Name:     env["DB_NAME"],
		},
		AuthorizationDADATA: config.AuthorizationDADATA{
			ApiKeyValue:    env["DADATA_API_KEY"],
			SecretKeyValue: env["DADATA_SECRET_KEY"],
		},
	}

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

	postgresDB := storage.NewGeoRepositoryDB(db)

	redisClient := redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})

	defer redisClient.Close()

	cache := storage.NewGeoRedis(redisClient)
	storagDB := storage.NewGeoRepositoryProxy(*postgresDB, cache)
	sevisDAdata := servisgeo.NewDadataService(config.AuthorizationDADATA)
	gs.geoProvider = app.NewGeoProvider(storagDB, sevisDAdata)

	err = postgresDB.ConnectToDB()

	if err != nil {
		log.Printf("Error conect DB %s", err)
	}

	if err := rpc.Register(gs); err != nil {
		log.Printf("Error Registretions rpc %v", err)
		return err
	}

	listen, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Printf("Eroor Listen %v", err)
		return err
	}
	defer listen.Close()

	log.Printf("%s",listen)

	log.Println("RPC сервер запущен и прослушивает порт :1234")
	rpc.Accept(listen)

	return nil
}

func (gs *GeoService) AddressSearch(query RequestAddressSearch, reply *[]*Address) error {
    addresses, err := gs.geoProvider.AddressSearch(query.Query)
    if err != nil {
        log.Printf("Error AddressSearch: %v", err)
        return err
    }

	for _,adres := range addresses{
		*reply = append(*reply, &Address{
			GeoLat: adres.GeoLat,
			GeoLon: adres.GeoLon,
			Result: adres.Result,
		})

	}

    return nil
}


func (gs *GeoService) AddressGeoCode(geocode RequestAddressGeocode, reply *[]*Address) error {
		addresses, err := gs.geoProvider.GeoCode(geocode.Lat,geocode.Lng)
		if err != nil {
			log.Printf("Error AddressGeoCode: %v", err)
			return err
		}
		// Просто присваиваем новое значение reply через косвенное разыменование
		for _,adres := range addresses{
			*reply = append(*reply, &Address{
				GeoLat: adres.GeoLat,
				GeoLon: adres.GeoLon,
				Result: adres.Result,
			})
	
		}
	
		return nil
	}
