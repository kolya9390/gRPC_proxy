package handler_geo

import (
	"encoding/json"
	"log"
	"net/http"

	geo_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/geo"
)

type GeoHandler struct {
	client geo_client.GeoClient
}

func NewGeoHandler(client geo_client.GeoClient) *GeoHandler {
	return &GeoHandler{client: client}
}

func (g *GeoHandler) SearchAPI(w http.ResponseWriter, r *http.Request) {

	var requestBody RequestAddressSearch

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("Decoder Body")
		return
	}

	addresses := g.client.SearchGeoAdres(geo_client.RequestAddressSearch(requestBody))

	// Конвертируйте объект ResponseAddress в JSON
	jsonResponse, err := json.Marshal(addresses)
	if err != nil {
		log.Println("Error marshalling JSON:", err)
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Error writing JSON response:", err)
		return
	}
}

func (g *GeoHandler) GeocodeAPI(w http.ResponseWriter, r *http.Request) {

	var requestBody RequestAddressGeocode

			if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
				log.Println("Decoder Body")
				return
			}

			addresses := g.client.GeoCoder(geo_client.RequestAddressGeocode(requestBody))

			// Конвертируйте объект ResponseAddress в JSON
			jsonResponse, err := json.Marshal(addresses)
			if err != nil {
				log.Println("Error marshalling JSON:", err)
				http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			_, err = w.Write(jsonResponse)
			if err != nil {
				log.Println("Error writing JSON response:", err)
				return
			}

}
