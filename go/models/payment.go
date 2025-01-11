package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Payment struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OrderID         string            `bson:"orderId" json:"orderId"`
	Course          primitive.ObjectID `bson:"course" json:"course"`
	User            primitive.ObjectID `bson:"user" json:"user"`
	Instructor      primitive.ObjectID `bson:"instructor" json:"instructor"`
	Amount          float64           `bson:"amount" json:"amount"`
	Currency        string            `bson:"currency" json:"currency"`
	Status          string            `bson:"status" json:"status"`
	PaymentMethod   string            `bson:"paymentMethod" json:"paymentMethod"`
	MidtransResponse interface{}      `bson:"midtransResponse" json:"midtransResponse"`
	PaymentLink     string            `bson:"paymentLink" json:"paymentLink"`
	CreatedAt       time.Time         `bson:"createdAt" json:"createdAt"`
	UpdatedAt       time.Time         `bson:"updatedAt" json:"updatedAt"`
}