package repo

import (
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/database"
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
	sqlc *database.Queries
}

func (up *UserRepository) GetUserByEmail(email string) bool {
	// row := global.MDB.Table(TableNameGoCrmUser).Where("usr_email = ?", email).First(&model.GoCrmUser{}).RowsAffected
	// return row != NumberNull
	user, err := up.sqlc.GetUserByEmail(ctx, email)
	if err != nil {
		return false
	}
	return user.UsrID != 0
}

func NewUserRepository() IUserRepository {
	return &UserRepository{
		sqlc: database.New(global.MDBC),
	}
}
