package routes

import (
	"github.com/amrilsyaifa/go_mongodb_rest_api/controllers"
	"github.com/amrilsyaifa/go_mongodb_rest_api/middleware"
	"github.com/amrilsyaifa/go_mongodb_rest_api/services"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func(authController *AuthRouteController) AuthRoute(routeGin *gin.RouterGroup, userService services.UserService) {
	router := routeGin.Group("/auth")

	router.POST("/register", authController.authController.SignUpUser)
	router.POST("/login", authController.authController.SignInUser)
	router.GET("/refresh", authController.authController.RefreshAccessToken)
	router.GET("/logout",middleware.DeserializeUser(userService), authController.authController.LogoutUser)
}