package app

import "context"

type UserProvider interface {
	GetUserIDs(user_id int64) (User, error)
	GetAllUser() ([]User, error)
	RegisterNewUser(ctx context.Context, name, email string, password string) (int64, error)
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password []byte `json:"password"`
}
