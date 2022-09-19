package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/amrilsyaifa/go_mongodb_rest_api/models"
	"github.com/amrilsyaifa/go_mongodb_rest_api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthServiceImplementation struct {
	collection 	*mongo.Collection
	ctx 		context.Context
}

func NewAuthServiceImplementation(collection *mongo.Collection, ctx context.Context) AuthService {
	return &AuthServiceImplementation{collection, ctx}
}

func (authService *AuthServiceImplementation) SignUpUser(user *models.SignUpInput) (*models.DBResponse, error) {
	/**
		mapping data user before insert to DB
	**/
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.PasswordConfirm = ""
	user.Verified = true 
	user.Role = "user"

	hashedPassword, _ := utils.HashPasword(user.Password)
	user.Password = hashedPassword

	// added the new user to the database with the InsertOne() function
	res, err := authService.collection.InsertOne(authService.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
	}

	// Create a unique index for the email field to ensure that no two users can have the same email address
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := authService.collection.Indexes().CreateOne(authService.ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	// I used the FindOne() function to find and return the user that was added to the database.
	var newUser *models.DBResponse
	query := bson.M{"_id": res.InsertedID}

	err = authService.collection.FindOne(authService.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (uc *AuthServiceImplementation) SignInUser(*models.SignInInput) (*models.DBResponse, error) {
	return nil, nil
}