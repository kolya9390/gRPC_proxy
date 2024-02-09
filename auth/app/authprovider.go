package app

import (
	passhesher "auth/infrastructure/tools/passHesher"
	"auth/storage"
	"context"
	"time"
)

type AuthApp struct {
	storageDB storage.UserRepo
	tokenTTL  time.Duration
}

func NewAuthApp(storage storage.UserRepo, tokenTTL time.Duration) *AuthApp {

	return &AuthApp{storageDB: storage,tokenTTL: tokenTTL}
}


func (a *AuthApp) RegisterNewUser(ctx context.Context,email string,password string) (int64, error){

	pass_hesh, err := passhesher.HashPassword(password)

	if err != nil {
		return 0 ,err
	}

	id,_ := a.storageDB.AddUser(ctx,email,[]byte(pass_hesh))

	return id,nil
	
}

func (a *AuthApp) Login(ctx context.Context,email string, password string) (string, error){
	panic("LOGIN")
}