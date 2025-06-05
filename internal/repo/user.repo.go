package repo

import (
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/model"
)

// type UserRepo struct{}

// func NewUserRepo() *UserRepo {
// 	return &UserRepo{}
// }

// func (ur *UserRepo) GetInfoUser() string {
// 	return "thientriel"
// }

// Interface version of UserRepo

type IUserRepository interface {
	GetUserByEmail(email string) bool
}

type UserRepository struct {
}

func (*UserRepository) GetUserByEmail(email string) bool {
	row := global.MDB.Table(TableNameGoCrmUser).Where("usr_email = ?", email).First(&model.GoCrmUser{}).RowsAffected
	return row != NumberNull
}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}
