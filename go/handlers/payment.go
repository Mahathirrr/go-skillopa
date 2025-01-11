package handlers

import (
	"context"
	"learnlit/database"
	"learnlit/models"
	"learnlit/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	var course models.Course
	err := database.DB.Collection("courses").FindOne(context.Background(), bson.M{"_id": courseObjID}).Decode(&course)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	var user models.User
	err = database.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": userObjID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	orderID := "ORDER-" + uuid.New().String()
	amount := int64(course.Price) // Convert float64 to int64

	resp, err := utils.CreatePaymentToken(orderID, amount, user, course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment token"})
		return
	}

	payment := models.Payment{
		OrderID:     orderID,
		Course:      courseObjID,
		User:        userObjID,
		Instructor:  course.Instructors[0],
		Amount:      course.Price,
		Currency:    course.Currency,
		Status:      "pending",
		PaymentLink: resp.RedirectURL,
	}

	_, err = database.DB.Collection("payments").InsertOne(context.Background(), payment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orderId":     orderID,
		"paymentLink": resp.RedirectURL,
		"token":       resp.Token,
	})
}

func GetPaymentStatus(c *gin.Context) {
	orderID := c.Param("orderId")

	var payment models.Payment
	err := database.DB.Collection("payments").FindOne(context.Background(), bson.M{"orderId": orderID}).Decode(&payment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	c.JSON(http.StatusOK, payment)
}

func HandlePaymentNotification(c *gin.Context) {
	var notification struct {
		OrderID           string `json:"order_id"`
		TransactionStatus string `json:"transaction_status"`
		FraudStatus       string `json:"fraud_status"`
	}

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payment models.Payment
	err := database.DB.Collection("payments").FindOne(context.Background(), bson.M{"orderId": notification.OrderID}).Decode(&payment)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	// Update payment status based on notification
	var status string
	switch notification.TransactionStatus {
	case "capture":
		if notification.FraudStatus == "challenge" {
			status = "pending"
		} else if notification.FraudStatus == "accept" {
			status = "success"
		}
	case "settlement":
		status = "success"
	case "cancel", "deny", "expire":
		status = "failed"
	case "pending":
		status = "pending"
	}

	_, err = database.DB.Collection("payments").UpdateOne(
		context.Background(),
		bson.M{"orderId": notification.OrderID},
		bson.M{"$set": bson.M{"status": status}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment status"})
		return
	}

	// If payment successful, enroll user in course
	if status == "success" {
		// Update user's enrolled courses
		_, err = database.DB.Collection("users").UpdateOne(
			context.Background(),
			bson.M{"_id": payment.User},
			bson.M{"$addToSet": bson.M{
				"enrolledCourses": bson.M{
					"course":     payment.Course,
					"enrolledOn": time.Now(),
				},
			}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enroll user"})
			return
		}

		// Update course enrollments
		_, err = database.DB.Collection("courses").UpdateOne(
			context.Background(),
			bson.M{"_id": payment.Course},
			bson.M{"$addToSet": bson.M{
				"meta.enrollments": bson.M{"id": payment.User},
			}},
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course enrollments"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

