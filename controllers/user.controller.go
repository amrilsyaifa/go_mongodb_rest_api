package controllers

import (
	"net/http"

	"github.com/amrilsyaifa/go_mongodb_rest_api/models"
	"github.com/amrilsyaifa/go_mongodb_rest_api/services"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(userService services.UserService) UserController {
	return UserController{userService}
}

func(userController *UserController) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}