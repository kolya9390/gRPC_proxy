package user_client

import "context"

type UserGetter interface {
	GetUserIDs(ctx context.Context, user_id int64) (User, error)
	GetListUsers(ctx context.Context) ([]User, error)
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
