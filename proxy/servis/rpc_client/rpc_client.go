package rpcclient

import (
	"context"
	"log"
	"net/rpc"
	"net/rpc/jsonrpc"

	"github.com/joho/godotenv"
	geo_provider "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type rpcGeoClient struct {
	client *rpc.Client
}

type jsonrpcGeoClient struct {
	client *rpc.Client
}

type grpcGeoClient struct {
	client *grpc.ClientConn
}

func NewGeoClient() GeoClient {

	env, err := godotenv.Read("client_app/.env")
	if err != nil {
		log.Fatal(err)
	}
	typeRPC := env["RPC_PROTOCOL"]
	switch typeRPC {
	case "rpc":
		client, err := rpc.Dial("tcp", "myNetwork:1234")
		if err != nil {
			log.Fatal("type-rpc", err)
		}
		return &rpcGeoClient{client: client}

	case "json-rpc":
		client, err := jsonrpc.Dial("tcp", "myNetwork:1234")
		if err != nil {
			log.Fatal("type-json", err)
		}
		return &jsonrpcGeoClient{client: client}

	case "grpc":
		client, err := grpc.Dial("server_rpc:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("did not connect %s", err)
		}
		return &grpcGeoClient{client: client}

	default:
		log.Fatalf("Unknown protocol: %s", typeRPC)
		return nil
	}
}

type Address struct {
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`
	Result string `json:"result"`
}

func (rgc *rpcGeoClient) SearchGeoAdres(query RequestAddressSearch) []Address {

	var result []Address
	err := rgc.client.Call("GeoService.AddressSearchRPC", query, &result)

	if err != nil {
		log.Println(err)
	}

	log.Println("Результат поиска:", result)

	return result
}

func (rgc *rpcGeoClient) GeoCoder(geocode RequestAddressGeocode) []Address {

	var result []Address
	err := rgc.client.Call("GeoService.AddressGeoCodeRPC", geocode, &result)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Результат поиска:", result)

	return result
}

func (c *jsonrpcGeoClient) SearchGeoAdres(query RequestAddressSearch) []Address {
	var result []Address
	err := c.client.Call("GeoService.AddressSearchRPC", query, &result)
	if err != nil {
		log.Println(err)
	}
	log.Println("Результат поиска:", result)
	return result
}

func (c *jsonrpcGeoClient) GeoCoder(geocode RequestAddressGeocode) []Address {
	var result []Address
	err := c.client.Call("GeoService.AddressGeoCodeRPC", geocode, &result)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Результат поиска:", result)
	return result
}

func (c *grpcGeoClient) SearchGeoAdres(query RequestAddressSearch) []Address {
	client := geo_provider.NewGeoProviderGRPCClient(c.client)
	// Вызываем метод SearchGeoAdres
	resp, err := client.AddressSearchGRPC(context.Background(), &geo_provider.RequestAddressSearch{
		Query: query.Query,
	})
	if err != nil {
		log.Printf("Error searching for address: %v", err)
		return nil
	}
	// Формируем ответ
	return []Address{
		{
			GeoLat: resp.Geolat,
			GeoLon: resp.Geolon,
			Result: resp.Result,
		},
	}
}

func (c *grpcGeoClient) GeoCoder(geocode RequestAddressGeocode) []Address {
	// Определяем gRPC клиент
	client := geo_provider.NewGeoProviderGRPCClient(c.client)
	// Вызываем метод GeoCoder
	resp, err := client.AddressGeoCodeGRPC(context.Background(), &geo_provider.RequestAddressGeocode{
		Lat: geocode.Lat,
		Lng: geocode.Lng,
	})
	if err != nil {
		log.Printf("Error geocoding address: %v", err)
		return nil
	}
	// Формируем ответ
	return []Address{
		{
			GeoLat: resp.Geolat,
			GeoLon: resp.Geolon,
			Result: resp.Result,
		},
	}
}
