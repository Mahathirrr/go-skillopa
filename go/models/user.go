package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name            string              `bson:"name" json:"name"`
	Email           string              `bson:"email" json:"email"`
	Password        string              `bson:"password,omitempty" json:"-"`
	Avatar          string              `bson:"avatar" json:"avatar"`
	Role            []string            `bson:"role" json:"role"`
	AuthProvider    string              `bson:"authProvider" json:"authProvider"`
	Cart            []primitive.ObjectID `bson:"cart" json:"cart"`
	Wishlist        []primitive.ObjectID `bson:"wishlist" json:"wishlist"`
	EnrolledCourses []EnrolledCourse    `bson:"enrolledCourses" json:"enrolledCourses"`
	Interests       []string            `bson:"interests" json:"interests"`
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type EnrolledCourse struct {
	Course     primitive.ObjectID `bson:"course" json:"course"`
	EnrolledOn time.Time         `bson:"enrolledOn" json:"enrolledOn"`
}