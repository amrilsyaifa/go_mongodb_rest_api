package controllers

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/amrilsyaifa/go_mongodb_rest_api/config"
	"github.com/amrilsyaifa/go_mongodb_rest_api/models"
	"github.com/amrilsyaifa/go_mongodb_rest_api/services"
	"github.com/amrilsyaifa/go_mongodb_rest_api/utils"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthController struct {
	authService	services.AuthService
	userService services.UserService
	ctx         context.Context
	collection  *mongo.Collection
	temp        *template.Template
}

func NewAuthController(authService services.AuthService, userService services.UserService, ctx context.Context, collection *mongo.Collection, temp *template.Template) AuthController {
	return AuthController{authService, userService, ctx, collection, temp}
}

func (authController *AuthController) SignUpUser(ctx *gin.Context) {
	// I validated the user’s input against the SignUpInput struct and returned an error if any of the rules were not satisfied.
	var user *models.SignUpInput

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// Next, I checked if the Password and PasswordConfirm values are equal.
	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "password do not match"})
		return
	}

	// Then I called the SignUpUser service with the user pointer to add the new user to the database.
	newUser, err := authController.authService.SignUpUser(user)
	if err != nil {
		if strings.Contains(err.Error(), "email already exist") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Generate Verification Code
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update User in Database
	authController.userService.UpdateUserById(newUser.ID.Hex(), "verificationCode", verificationCode)

	var firstName = newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// 👇 Send Email
	emailData := utils.EmailData{
		URL:       config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, authController.temp, "verificationCode.html")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "There was an error sending email"})
		return
	}

	message := "We sent an email with a verification code to " + user.Email

	// Lastly, I sent a JSON response to the user assuming there was no error. In an upcoming tutorial, we’ll send an email verification code to the user’s email.
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (authController *AuthController) SignInUser(ctx *gin.Context) {
	var credentials *models.SignInInput

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	user, err := authController.userService.FindUserByEmail(credentials.Email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or password"})
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	if err := utils.VerifyPassword(user.Password, credentials.Password); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate Tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivateKey)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refresh_token, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (authController *AuthController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	cookie, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	sub, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := authController.userService.FindUserByID(fmt.Sprint(sub))
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivateKey )
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", access_token, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token})
}

func (authController *AuthController) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (ac *AuthController) VerifyEmail(ctx *gin.Context) {

	code := ctx.Params.ByName("verificationCode")
	verificationCode := utils.Encode(code)

	query := bson.D{{Key: "verificationCode", Value: verificationCode}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "verified", Value: true}}}, {Key: "$unset", Value: bson.D{{Key: "verificationCode", Value: ""}}}}
	result, err := ac.collection.UpdateOne(ac.ctx, query, update)
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "success", "message": "Could not verify email address"})
		return
	}

	fmt.Println(result)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Email verified successfully"})

}