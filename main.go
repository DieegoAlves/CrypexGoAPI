package main

import (
	"github.com/DieegoAlves/CrypexGoAPI/src/controller"
	"github.com/DieegoAlves/CrypexGoAPI/src/repositories"
	"github.com/DieegoAlves/CrypexGoAPI/src/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"os"
)

func main() {
	r := gin.Default()
	db := &gorm.DB{
		Config:       nil,
		Error:        nil,
		RowsAffected: 0,
		Statement:    nil,
	}

	userRepository := repositories.NewUserRepository(db)

	userService := services.NewUserService(userRepository)

	userController := controller.NewUserController(userService)

	r.POST("/user", userController.CreateUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := r.Run(":" + port)
	if err != nil {
		panic(err)
	}
}
