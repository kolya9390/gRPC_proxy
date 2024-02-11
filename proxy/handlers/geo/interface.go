package handler_geo

import "net/http"

type GeoServiceController interface {
	SearchAPI(w http.ResponseWriter, r *http.Request)
	GeocodeAPI(w http.ResponseWriter, r *http.Request)
}

type Address struct {
	GeoLat string `json:"lat"`
	GeoLon string `json:"lon"`

	Result string `json:"result"`
}

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type ResponseAddress struct {
	Addresses []Address `json:"addresses"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}
