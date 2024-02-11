package handler_auth

import "net/http"

type AuthServiceController interface {
	Login(w http.ResponseWriter, r *http.Request)
	Registeretion(w http.ResponseWriter, r *http.Request)
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
/*
type Token struct {
	Token string `json:"Bearer"`
}
*/

type ReqpLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RespLogin struct {
	Token string `json:"Authorization"`
}

type ReqRegister struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RespRegister struct {
	Statuse string `json:"status"`
	User_ID int64  `json:"user_id"`
}