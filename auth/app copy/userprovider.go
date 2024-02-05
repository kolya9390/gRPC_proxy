package app

import (
	"user/storage"
)

type UserService struct {
	storeg storage.UserRepository

}

func NewGeoProvider(storagDB storage.UserRepository)*UserService{
	return &UserService{storeg: storagDB}
}

func (us *UserService) GetUserIDs(id string) (User, error) {

    user, err := us.storeg.GetUserIDs(id)

    if err != nil {
        return User{},err
    }

    return User{user.Email,user.Password},nil
}

func (us *UserService) GetAllUser() ([]User, error) {

    var result []User

    users, err := us.storeg.GetUsers()

    if err != nil {
        return nil,err
    }

    for _, user := range users{
        result = append(result, User{user.Email,user.Password})
    }

    return result,nil
}
    