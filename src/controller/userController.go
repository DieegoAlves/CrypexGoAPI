package controller

import (
	"github.com/DieegoAlves/CrypexGoAPI/src/entities"
	"github.com/DieegoAlves/CrypexGoAPI/src/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"time"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type UserController struct {
	service services.UserService
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUserController(userService services.UserService) UserController {
	return UserController{
		service: userService,
	}
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	user := entities.User{}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	user.ID = uuid.New()
	err := ctx.ShouldBindBodyWithJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = user.VerifyFields()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = u.service.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (u *UserController) Login(ctx *gin.Context) {

	var userCredentials credentials
	err := ctx.ShouldBindJSON(&userCredentials)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	if userCredentials.Password == "" || userCredentials.Username == "" {
		ctx.JSON(http.StatusBadRequest, "empty field")
	}

	valid, err := u.service.VerifyCredentials(userCredentials.Username, userCredentials.Password)
	if err != nil || !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := u.service.GenerateJWT(userCredentials.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
	})
}

func (u *UserController) GetUser(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, err := u.service.FindByUsername(username.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"username": user.Username,
	})
}

func (u *UserController) UpdateUsername(ctx *gin.Context) {
	currentUsername, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	newUsername := ctx.PostForm("new_username")

	if err := u.service.UpdateUsername(currentUsername.(string), newUsername); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Username update successfully. Log in again!"})
}

func (u *UserController) UpdateBio(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
	newBio := ctx.PostForm("new_bio")

	if err := u.service.UpdateBio(username.(string), newBio); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Bio update successfully"})
}

func (u *UserController) DeleteUser(ctx *gin.Context) {
	var userCredentials credentials
	err := ctx.ShouldBindJSON(&userCredentials)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}

	nameCheck, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userCredentials.Username = nameCheck.(string)

	if userCredentials.Password == "" || userCredentials.Username == "" {
		ctx.JSON(http.StatusUnauthorized, "empty field")
		return
	}

	valid, err := u.service.VerifyCredentials(userCredentials.Username, userCredentials.Password)
	if err != nil || !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	user, err := u.service.FindByUsername(userCredentials.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = u.service.DeleteUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
