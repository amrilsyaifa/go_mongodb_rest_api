package routes

import (
	"github.com/amrilsyaifa/go_mongodb_rest_api/controllers"
	"github.com/amrilsyaifa/go_mongodb_rest_api/middleware"
	"github.com/amrilsyaifa/go_mongodb_rest_api/services"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewRouteUserController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}	
}

func (userController *UserRouteController) UserRoute(routeGin *gin.RouterGroup, userService services.UserService) {
	router := routeGin.Group("users")

	router.Use(middleware.DeserializeUser(userService))
	router.GET("/me", userController.userController.GetMe)
}