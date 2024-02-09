package storage

import (
	"context"
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewAuthRepositoryDB(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (s *UserRepo) AddUser(ctx context.Context, email string, passHash []byte) (int64, error) {

	var id int64
	err := s.db.QueryRowx("INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id", email, passHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (s *UserRepo) GetUser(ctx context.Context, email string) (User, error) {

	var user User
	err := s.db.Get(&user, `SELECT email, password FROM users WHERE email = $1`, email)

	if err != nil {
		return User{}, err
	}

	return user, nil
}
