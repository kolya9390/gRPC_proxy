package storage

type UserRepository interface {
	GetUserIDs(userID string) (User,error) // Вставка в DB's
	GetUsers() ([]User, error) // Получаем данные из базы
}

type User struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
}