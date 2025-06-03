//go:build wireinject

package wire

import (
	"github.com/google/wire"
	"go-ecommerce-backend-api.com/internal/controller"
	"go-ecommerce-backend-api.com/internal/repo"
	"go-ecommerce-backend-api.com/internal/service"
)

func InitUserRouterHandler() (*controller.UserController,error){
	wire.Build(
		repo.NewUserRepository,
		service.NewUserService,
		controller.NewUserController,
	)
	return new(controller.UserController),nil
}
