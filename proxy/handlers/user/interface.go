package handler_user

import "net/http"

type UserServiceController interface {
	GetUserID(w http.ResponseWriter, r *http.Request)
	GetUsers(w http.ResponseWriter, r *http.Request)
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}

type RequestGetUserID struct{
	User_ID int64 `json:"user_id"`
}

type RespGetUserID struct {
	User User	`json:"user"`
}

type ReqGetUsers struct {

}

type RespGetUsers struct {
	Users []User `json:"List_Users"`
}