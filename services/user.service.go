package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/amrilsyaifa/go_mongodb_rest_api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


type UserServiceImplementation struct {
	collection 	*mongo.Collection
	ctx 		context.Context
}

func NewUserServiceImplementation(collection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImplementation{collection, ctx}
}

func (usrService *UserServiceImplementation) FindUserByID(id string) (*models.DBResponse, error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	
	var user *models.DBResponse

	query := bson.M{"_id": oid}
	err := usrService.collection.FindOne(usrService.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments{
			return &models.DBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (usrService *UserServiceImplementation) FindUserByEmail(email string) (*models.DBResponse, error) {
	var user *models.DBResponse

	query := bson.M{"email": strings.ToLower(email)}
	err := usrService.collection.FindOne(usrService.ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.DBResponse{}, err
		}

		return nil, err
	}

	return user, nil
}

func (uc *UserServiceImplementation) UpdateUserById(id string, field string, value string) (*models.DBResponse, error) {
	userId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: userId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	fmt.Print(result.ModifiedCount)
	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}

func (uc *UserServiceImplementation) UpdateOne(field string, value interface{}) (*models.DBResponse, error) {
	query := bson.D{{Key: field, Value: value}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: field, Value: value}}}}
	result, err := uc.collection.UpdateOne(uc.ctx, query, update)

	fmt.Print(result.ModifiedCount)
	if err != nil {
		fmt.Print(err)
		return &models.DBResponse{}, err
	}

	return &models.DBResponse{}, nil
}