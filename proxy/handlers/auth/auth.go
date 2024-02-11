package handler_auth

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	auth_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/auth"
)

type AuthHandler struct {
	client auth_client.Auther
}

func NewAuthHandler(client auth_client.Auther) *AuthHandler {
	return &AuthHandler{client: client}
}

func (a *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	
	var requestBody ReqpLogin

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("Decoder Body")
		return
	}

	token, err := a.client.Login(context.Background(),requestBody.Email,requestBody.Password)

	if err != nil {
		log.Println("a.client.Login")
	}

	resp := RespLogin{Token: token}
	jsonResponse, err := json.Marshal(resp)
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

func (a *AuthHandler) Registeretion(w http.ResponseWriter, r *http.Request) {

	var requestBody ReqRegister

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("Decoder Body")
		return
	}

	user_id, err := a.client.Registeretion(context.Background(),
	requestBody.Name,requestBody.Email,requestBody.Password)

	resp := RespRegister{User_ID: user_id,Statuse: "Successfully registrations"}

	if err != nil {
		log.Println("a.client.Registeretion")
		resp = RespRegister{User_ID: user_id,Statuse: "unsuccessfully registrations"}
	}

	if user_id == 0 {
		resp = RespRegister{User_ID: user_id,Statuse: "unsuccessfully registrations"}
	}

	jsonResponse, err := json.Marshal(resp)
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