package models

import (
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CustomClaims -> a claimed-token for jwt pkg
type CustomClaims struct {
	Username 	string `json:"username"`
	Email		string `json:"email"`
	jwt.StandardClaims
}

// RequiredDataToClaims -> is the struct with the required data that is needed for claiming the token
type RequiredDataToClaims struct {
	Username	string
	Email		string
	ID 			primitive.ObjectID
}