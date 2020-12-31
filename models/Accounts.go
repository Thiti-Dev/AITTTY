package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Accounts -> An struct that describe the user-account key field
type Accounts struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username  string             `json:"username,omitempty" bson:"username,omitempty" binding:"required" validate:"required"`
	Email     string             `json:"email,omitempty" bson:"email,omitempty" binding:"required" validate:"required,email"`
	CreatedAt time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty" binding:"required"`
	UpdatedAt time.Time          `json:"updatedAt,omitempty" bson:"updatedAt,omitempty" binding:"required"`
}