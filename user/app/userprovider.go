package app

import (
	"context"
	passhesher "user/infrastructure/tools/passHesher"
	"user/storage"
)

type UserService struct {
	storege storage.UserRepository

}

func NewUserProvider(storagDB storage.UserRepository)*UserService{
	return &UserService{storege: storagDB}
}

func (us *UserService) GetUserIDs(id int64) (User, error) {

    user, err := us.storege.GetUserIDs(id)

    if err != nil {
        return User{},err
    }

    return User{user.ID,user.Name,user.Email,user.Password},nil
}

func (us *UserService) GetAllUser() ([]User, error) {

    var result []User

    users, err := us.storege.GetUsers()

    if err != nil {
        return nil,err
    }

    for _, user := range users{
        result = append(result, User{user.ID,user.Name,user.Email,user.Password})
    }

    return result,nil
}
    
func (a *UserService) RegisterNewUser(ctx context.Context,name,email string,password string) (int64, error){

	pass_hesh, err := passhesher.HashPassword(password)

	if err != nil {
		return 999999 ,err
	}

	id,_ := a.storege.AddUser(ctx,name,email,[]byte(pass_hesh))

	return id,nil
	
}