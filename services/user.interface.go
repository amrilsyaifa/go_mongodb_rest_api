package services

import "github.com/amrilsyaifa/go_mongodb_rest_api/models"

type UserService interface {
	FindUserByID(string)(*models.DBResponse, error)
	FindUserByEmail(string)(*models.DBResponse, error)
	UpdateUserById(id string, field string, value string) (*models.DBResponse, error)
	UpdateOne(field string, value interface{}) (*models.DBResponse, error)
}