package handler_user

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	user_client "github.com/kolya9390/gRPC_GeoProvider/client_Proxy/servis/rpc_client/user"
)

type UserHandler struct {
	client user_client.UserGetter
}

func NewUserHandler(client user_client.UserGetter) *UserHandler {
	return &UserHandler{client: client}
}

func (u *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	var requestBody RequestGetUserID

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		log.Println("Decoder Body")
		return
	}

	user,err := u.client.GetUserIDs(context.Background(),requestBody.User_ID)

	if err != nil {
		log.Println("u.client.GetUserIDs")
		return
	}

	jsonResponse, err := json.Marshal(user)
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

func (u *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request){

	users,err := u.client.GetListUsers(context.Background())

	if err != nil {
		log.Println("u.client.GetListUsers")
		return
	}

	jsonResponse, err := json.Marshal(users)
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