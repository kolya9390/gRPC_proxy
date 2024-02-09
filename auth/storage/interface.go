package storage

import "context"

type AuthRepository interface {
	// Вставка в DB's
	AddUser(ctx context.Context, email string, passHash []byte) (int64, error)
	GetUser(ctx context.Context, email string) (User, error)
}

type User struct {
	Email    string
	PassHash []byte
}