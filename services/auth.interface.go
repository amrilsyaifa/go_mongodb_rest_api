package services

import "github.com/amrilsyaifa/go_mongodb_rest_api/models"

/*
	Go language interfaces are different from other languages. In Go language, the interface is a custom type that is used to specify a set of one or more method signatures and the interface is abstract, so you are not allowed to create an instance of the interface.
*/

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBResponse, error)
	SignInUser(*models.SignInInput) (*models.DBResponse, error)
}
