package service

import (
	"go-ecommerce-backend-api.com/internal/repo"
	"go-ecommerce-backend-api.com/pkg/response"
)

// import "go-ecommerce-backend-api.com/internal/repo"

// type UserService struct {
// 	userRepo *repo.UserRepo
// }

// func NewUserService() *UserService {
// 	return &UserService{
// 		userRepo: repo.NewUserRepo(),
// 	}
// }
// func (us *UserService) GetInfoUser() string {
// 	return us.userRepo.GetInfoUser()
// }

// interface UserServiceInterface

type IUserService interface {
	Regisger(email string, purpose string) int
}

type userService struct {
	userRepo repo.IUserRepository
	//..
}

func NewUserService(
	userRepo repo.IUserRepository,
) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

// Regisger implements IUserService.
func (us *userService) Regisger(email string, purpose string) int {
	// TODO: implement registration logic

	if us.userRepo.GetUserByEmail(email) {
		return response.ErrcodeUserHasExist
	}

	return response.ErrCodeSuccess
}
