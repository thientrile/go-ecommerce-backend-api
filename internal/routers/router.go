package routers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go-ecommerce-backend-api.com/internal/middlewares"
)

func AA() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("before AA")
		c.Next() // Call the next handler in the chain
		fmt.Println("after AA")
	}
}

func BB() gin.HandlerFunc {
	return func(c *gin.Context) {

		fmt.Println("before BB")
		c.Next() // Call the next handler in the chain
		fmt.Println("after BB")
	}
}

func CC(c *gin.Context) {
	fmt.Println("before CC")
	c.Next() // Call the next handler in the chain
	fmt.Println("after CC")
}
func NewRouter() *gin.Engine {
	r := gin.Default()
	// use the middleware
	r.Use(middlewares.AuthMiddleware(), AA(), BB(), CC)
	// v1 := r.Group("/v1/2025") //v1 2025
	// {
	// 	v1.GET("/user/:id", c.NewUserController().GetUserById)
	// 	// v1.PATCH("/ping", Pong)
	// 	// v1.PUT("/ping", Pong)
	// 	// v1.POST("/ping", Pong)
	// 	// v1.DELETE("/ping", Pong)
	// 	// v1.HEAD("/ping", Pong)
	// 	// v1.OPTIONS("/ping", Pong)
	// }

	v2 := r.Group("/v2/2025") //v2 2025
	{
		v2.GET("/ping", Pong)
		v2.PATCH("/ping", Pong)
		v2.PUT("/ping", Pong)
		v2.POST("/ping", Pong)
		v2.DELETE("/ping", Pong)
		v2.HEAD("/ping", Pong)
		v2.OPTIONS("/ping", Pong)
	}
	return r
}

func Pong(c *gin.Context) {
	name := c.Param("name")
	uid := c.Query("uid")
	c.JSON(http.StatusOK, gin.H{
		"message": "pong" + name + "uid:" + uid,
		"code":    http.StatusOK,
		"users":   []string{"user1", "user2", "user3"},
	})
}
