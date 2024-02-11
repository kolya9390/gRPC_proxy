package geo_client

type GeoClient interface {
	SearchGeoAdres(query RequestAddressSearch) []Address
	GeoCoder(geocode RequestAddressGeocode) []Address
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}