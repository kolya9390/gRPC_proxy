package user

type UserService interface {
	GetUserIDs(user_id string)
	AddUser()
}