package repo

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
	// TODO: implement actual logic
	return true
}

func NewUserRepository() IUserRepository {
	return &UserRepository{}
}
