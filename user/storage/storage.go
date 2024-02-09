package storage

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepositoryDB(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Создание таблицы пользователей
func (d *UserRepo) ConnectToDB() error {
	sqlStatementUser := `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
	name TEXT,
    email TEXT,
    pass_hash TEXT
);`

	_, err := d.db.Exec(sqlStatementUser)
	return err
}

// Получение пользователя по ID
func (s *UserRepo) GetUserIDs(userID int64) (User, error) {
	var user User
	err := s.db.Get(&user, `SELECT id, name, email FROM users WHERE id = $1`, userID)

	if err != nil {
		return User{},err
	}
	return user, nil
}

// Получение списка всех пользователей
func (s *UserRepo) GetUsers() ([]User, error) {
	var users []User
	err := s.db.Select(&users, `SELECT id, name, email FROM users`)
	if err != nil {
		return nil, err
	}
	return users, nil
}



func (s *UserRepo) AddUser(ctx context.Context,name, email string, passHash []byte) (int64, error) {

	var id int64

	var count int
    err := s.db.QueryRowx("SELECT COUNT(*) FROM users WHERE email = $1", email).Scan(&count)
    if err != nil {
        return 0, err
    }
    if count > 0 {
        // Пользователь с таким email уже существует, вернуть ошибку
        return 0, errors.New("user with this email already exists")
    }

	/// Если пользовотель с такой почтой уже есть , то не надо добавлять его еще раз!! фикс
	err = s.db.QueryRowx("INSERT INTO users (name, email, pass_hash) VALUES ($1, $2, $3) RETURNING id",name, email, passHash).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

