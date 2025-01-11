package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Instructor struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name         string              `bson:"name" json:"name"`
	Slug         string              `bson:"slug" json:"slug"`
	Email        string              `bson:"email" json:"email"`
	Headline     string              `bson:"headline" json:"headline"`
	Bio          string              `bson:"bio" json:"bio"`
	Avatar       string              `bson:"avatar" json:"avatar"`
	Courses      []primitive.ObjectID `bson:"courses" json:"courses"`
	Social       Social              `bson:"social" json:"social"`
	BankAccount  BankAccount         `bson:"bankAccount" json:"bankAccount"`
	CreatedBy    primitive.ObjectID   `bson:"createdBy" json:"createdBy"`
	Meta         InstructorMeta      `bson:"meta" json:"meta"`
	CreatedAt    time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt    time.Time           `bson:"updatedAt" json:"updatedAt"`
}

type Social struct {
	Website  string `bson:"website" json:"website"`
	Twitter  string `bson:"twitter" json:"twitter"`
	Youtube  string `bson:"youtube" json:"youtube"`
	LinkedIn string `bson:"linkedin" json:"linkedin"`
	Facebook string `bson:"facebook" json:"facebook"`
}

type BankAccount struct {
	BankName         string `bson:"bankName" json:"bankName"`
	AccountNumber    string `bson:"accountNumber" json:"accountNumber"`
	AccountHolderName string `bson:"accountHolderName" json:"accountHolderName"`
}

type InstructorMeta struct {
	Enrollments   int     `bson:"enrollments" json:"enrollments"`
	TotalReviews  int     `bson:"totalReviews" json:"totalReviews"`
	TotalEarnings float64 `bson:"totalEarnings" json:"totalEarnings"`
}