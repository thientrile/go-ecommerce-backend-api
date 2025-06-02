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

type userRepository struct {
}

func (*userRepository) GetUserByEmail(email string) bool {
	// TODO: implement actual logic
	return true
}

func NewUserRepository() IUserRepository {
	return &userRepository{}
}
