package main

import (
	"github.com/DieegoAlves/CrypexGoAPI/src/controller"
	"github.com/DieegoAlves/CrypexGoAPI/src/database"
	"github.com/DieegoAlves/CrypexGoAPI/src/entities"
	"github.com/DieegoAlves/CrypexGoAPI/src/middlewares"
	"github.com/DieegoAlves/CrypexGoAPI/src/repositories"
	"github.com/DieegoAlves/CrypexGoAPI/src/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()
	db := database.ConnectDB()

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	err := db.AutoMigrate(entities.User{})
	if err != nil {
		log.Fatal(err)
		return
	}

	secure := r.Group("/secure")
	secure.Use(middlewares.JWTAuthMiddleware())
	{
		secure.GET("/auth", AuthHandler)

		secure.GET("/profile", userController.GetUser)

		secure.PUT("/update/username", userController.UpdateUsername)

		secure.PUT("/update/bio", userController.UpdateBio)

		secure.DELETE("/delete", userController.DeleteUser)
	}

	r.POST("/user", userController.CreateUser)

	r.POST("/login", userController.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
		return
	}
}

func AuthHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Welcome to the protected route!",
	})
}
