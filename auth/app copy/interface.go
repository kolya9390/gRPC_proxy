package app


type UserProvider interface {
	GetUserIDs(id string) (User, error)
	GetAllUser() ([]User, error)
}

type User struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
}