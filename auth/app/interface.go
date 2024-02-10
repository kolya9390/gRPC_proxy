package app

import "context"

type Auther interface {
	RegisterNewUser(ctx context.Context,name,email,password string) (int64, error)
	Login(ctx context.Context,email string, password string) (string, error)
}