package user

import "context"

type UserService interface {
	GetListUsers(ctx context.Context) ([]User,error)
	GetUserIDs(ctx context.Context, user_id int64) (User, error)
	Register(ctx context.Context,name, email, password string) (int64, error)
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
