package initialize

import (
	"go-ecommerce-backend-api.com/global"
	"go-ecommerce-backend-api.com/internal/database"
	"go-ecommerce-backend-api.com/internal/service"
	"go-ecommerce-backend-api.com/internal/service/impl"
)

func InitServiceInterface() {
	queries := database.New(global.MDBC)
	service.InitUserLogin(impl.NewUserLoginImpl(queries))
}
