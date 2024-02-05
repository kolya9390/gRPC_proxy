package app

import servicegeo "github.com/kolya9390/gRPC_GeoProvider/server_rpc/servis_geo"

type GeoProvider interface {
	AddressSearch(input string) ([]*servicegeo.Address, error)
	GeoCode(lat, lng string) ([]*servicegeo.Address, error)
}