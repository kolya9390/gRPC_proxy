package storage

import (
	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewGeoRepositoryDB(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// создание таблицы
func (d *UserRepo) ConnectToDB() error {

	sqlStatementUser := `
CREATE TABLE IF NOT EXISTS user (
    id SERIAL PRIMARY KEY,
    email text,
	password text

);`

	_, err := d.db.Exec(sqlStatementUser)
	if err != nil {
		return err
	}

	return nil

}

func (s *UserRepo) GetUserIDs(userID string) (User, error) {

	var user User

	err := s.db.Select(&user, `
	SELECT email as Email, password as Password 
	FROM user
	WHERE id  LIKE $1`, "%"+userID+"%")
	if err != nil {
		return User{}, err
	}

	return user, nil
}


func (s *UserRepo) GetUsers() ([]User,error){
	var user []User

	err := s.db.Select(&user,`
	SELECT email as Email, password as Password 
	FROM user`)
	if err != nil {
		return nil, err
	}

	return user,nil
}
