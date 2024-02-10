package app

import (
	jwt_token "auth/infrastructure/tools/jwt"
	passhesher "auth/infrastructure/tools/passHesher"
	"auth/service/user"
	"context"
	"errors"
	"fmt"
)

type AuthService struct {
	service       user.UserService
	mapsUsers_ids map[string]int64
}

func NewAuthProvider(service user.UserService) *AuthService {
	return &AuthService{service: service, mapsUsers_ids: make(map[string]int64)}
}

func (a *AuthService) Login(ctx context.Context, email string, password string) (string, error) {

	var user user.User

	var err error

	if user_id,OK := a.mapsUsers_ids[email]; !OK {
		user,err = a.getUsers(ctx,email)
		if err != nil{
			return "",fmt.Errorf("getUsers err: %s",err)
		}
	}else{
	user, err = a.service.GetUserIDs(ctx, user_id)
	if err != nil {
		return "", fmt.Errorf("GetUserIDs err: %s",err)
	}
	}


	token, err := jwt_token.NewToken(jwt_token.User{
        ID:       user.ID,
        Email:    user.Email,
    })
	if err != nil {
		return "Error", err
	}

	truePassword := passhesher.CheckPassword(string(user.Password), password)
	if !truePassword {
		return "incorrect password", errors.New("incorrect password")
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, name, email, password string) (int64, error) {

	user_id, err := a.service.Register(ctx, name, email, password)
	if err != nil {
		return 0, err
	}
	a.mapsUsers_ids[email]=user_id
	return user_id, nil
}

func (a *AuthService) getUsers(ctx context.Context,email string) (user.User,error){
	users,err :=  a.service.GetListUsers(ctx)
	if err != nil {
		return user.User{},err
	}

	for _,user := range users{
		if email == user.Email{
			//a.mapsUsers_ids[email]=user.ID
			return user,nil
		}
	}

	return user.User{}, fmt.Errorf("email not faund")
}