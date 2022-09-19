package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
	Whats is struct in go?
	In each of the tutorials, we made use of bson.A, bson.D, and bson.M, which represent arrays, documents, and maps. However, these are primitive data structures that are part of the MongoDB Go driver, and not necessarily the best way for interacting with data, both within the application and the database.
	We're going to look at an alternative way to interact with data through the MongoDB Go driver operations. This time we're going to map MongoDB document fields to native Go data structures.

	reference https://www.mongodb.com/blog/post/quick-start-golang--mongodb--modeling-documents-with-go-data-structures
*/

/*
	The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.
*/
// 👈 SignUpInput struct
type SignUpInput struct {
	Name				string		`json:"name" bson:"name" binding:"required"`
	Email				string		`json:"email" bson:"email" binding:"required"`
	Password			string		`json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm		string		`json:"passwordConfirm" bson:"passwordConfirm,omitempty" binding:"required"`
	Role 				string		`json:"role" bson:"role"`
	Verified			bool		`json:"verified" bson:"verified"`
	CreatedAt			time.Time	`json:"created_at" bson:"created_at"`
	UpdatedAt			time.Time	`json:"updated_at" bson:"updated_at"`
}

// 👈 SignInInput struct
type SignInInput struct {
	Email				string		`json:"email" bson:"email" binding:"required"`
	Password			string		`json:"password" bson:"password" binding:"required"`
}

/*
	DBResponse struct to define the fields that will be returned by MongoDB.
*/
// 👈 DBResponse struct
type DBResponse struct {
	ID					primitive.ObjectID	`json:"id" bson:"_id"`
	Name				string				`json:"name" bson:"name"`
	Email				string				`json:"email" bson:"email"`
	Password			string				`json:"password" bson:"password"`
	PasswordConfirm		string				`json:"passwordConfirm,omitempty" bson:"passwordConfirm,omitempty"`
	Role 				string				`json:"role" bson:"role"`
	Verified			bool				`json:"verified" bson:"verified"`
	CreatedAt			time.Time			`json:"created_at" bson:"created_at"`
	UpdatedAt			time.Time			`json:"updated_at" bson:"updated_at"`
}

// 👈 UserResponse struct
type UserResponse struct {
	ID        		primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      		string             `json:"name,omitempty" bson:"name,omitempty"`
	Email     		string             `json:"email,omitempty" bson:"email,omitempty"`
	Role      		string             `json:"role,omitempty" bson:"role,omitempty"`
	CreatedAt 		time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt 		time.Time          `json:"updated_at" bson:"updated_at"`	
}

/*
	a function to filter out the sensitive fields
*/
// 👈 FilteredResponse stru
func FilteredResponse(user *DBResponse) UserResponse {
	return UserResponse{
		ID:			user.ID,
		Email: 		user.Email,
		Name: 		user.Name,
		Role: 		user.Role,
		CreatedAt: 	user.CreatedAt,
		UpdatedAt: 	user.UpdatedAt,
	}
}