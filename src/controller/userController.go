package controller

import (
	"github.com/DieegoAlves/CrypexGoAPI/src/entities"
	"github.com/DieegoAlves/CrypexGoAPI/src/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type UserController struct {
	service services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{
		service: userService,
	}
}

func (u *UserController) CreateUser(ctx *gin.Context) {
	user := entities.User{}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

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

	ctx.JSON(http.StatusOK, user)
}
