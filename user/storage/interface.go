package storage

import "context"

type UserRepository interface {
	GetUserIDs(userID int64) (User, error) // 
	GetUsers() ([]User, error)              // Получаем данные из базы
	AddUser(ctx context.Context,name, email string, passHash []byte) (int64, error)//
	//GetUser(ctx context.Context, email string) (User, error)
}

type User struct {
    ID       int64  `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password []byte `json:"password"`
}
