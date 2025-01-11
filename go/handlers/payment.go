package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"learnlit/database"
	"learnlit/models"
	"learnlit/utils"
)

func CreatePayment(c *gin.Context) {
	var input struct {
		CourseID string `json:"courseId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))
	courseObjID, _ := primitive.ObjectIDFromHex(input.CourseID)

	// Get course details
	var course models.Course
	err := database.DB.Collection("courses").FindOne(context.Background(), 
		bson.M{"_id": courseObjID}).Decode(&course)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if course.Pricing != "Paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This course is free"})
		return
	}

	// Check if user already enrolled
	var user models.User
	err = database.DB.Collection("users").FindOne(context.Background(), 
		bson.M{
			"_id": userObjID,
			"enrolledCourses.course": courseObjID,
		}).Decode(&user)
	
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Already enrolled in this course"})
		return
	}

	// Create unique order ID
	orderID := fmt.Sprintf("ORDER-%s", uuid.New().String())

	// Create payment record
	payment := models.Payment{
		OrderID:    orderID,
		Course:     courseObjID,
		User:       userObjID,
		Instructor: course.Instructors[0],
		Amount:     course.Price,
		Currency:   course.Currency,
		Status:     "pending",
	}

	// Get payment token from Midtrans
	transaction, err := utils.CreatePaymentToken(orderID, course.Price, user, course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment token"})
		return
	}

	payment.PaymentLink = transaction.RedirectURL
	
	_, err = database.DB.Collection("payments").InsertOne(context.Background(), payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId": orderID,
		"paymentLink": transaction.RedirectURL,
		"token": transaction.Token,
	})
}