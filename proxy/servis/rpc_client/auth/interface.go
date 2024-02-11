package auth_client

import "context"

type Auther interface {
	Login(ctx context.Context,email,password string) (string,error)
	Registeretion(ctx context.Context,name,email,password string)(int64,error)
}