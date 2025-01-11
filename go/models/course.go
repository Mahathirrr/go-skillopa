package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Course struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string              `bson:"title" json:"title"`
	Subtitle    string              `bson:"subtitle" json:"subtitle"`
	Slug        string              `bson:"slug" json:"slug"`
	Description string              `bson:"description" json:"description"`
	Category    string              `bson:"category" json:"category"`
	SubCategory string              `bson:"subCategory" json:"subCategory"`
	Language    string              `bson:"language" json:"language"`
	Duration    string              `bson:"duration" json:"duration"`
	CoverImage  string              `bson:"coverImage" json:"coverImage"`
	PostedBy    primitive.ObjectID   `bson:"postedBy" json:"postedBy"`
	Instructors []primitive.ObjectID `bson:"instructors" json:"instructors"`
	Level       string              `bson:"level" json:"level"`
	Pricing     string              `bson:"pricing" json:"pricing"`
	Currency    string              `bson:"currency" json:"currency"`
	Price       float64             `bson:"price" json:"price"`
	Published   bool                `bson:"published" json:"published"`
	Meta        CourseMeta          `bson:"meta" json:"meta"`
	CreatedAt   time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type CourseMeta struct {
	Enrollments     []Enrollment `bson:"enrollments" json:"enrollments"`
	Rating          float64      `bson:"rating" json:"rating"`
	NumberOfRatings int         `bson:"numberOfRatings" json:"numberOfRatings"`
}

type Enrollment struct {
	ID         primitive.ObjectID `bson:"id" json:"id"`
	EnrolledOn time.Time         `bson:"enrolledOn" json:"enrolledOn"`
}