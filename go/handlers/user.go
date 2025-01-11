package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"learnlit/database"
	"learnlit/models"
	"learnlit/utils"
)

func AddToCart(c *gin.Context) {
	var input struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))
	courseID, _ := primitive.ObjectIDFromHex(input.ID)

	result, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$addToSet": bson.M{"cart": courseID}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course already in cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RemoveFromCart(c *gin.Context) {
	courseID := c.Param("id")
	courseObjID, _ := primitive.ObjectIDFromHex(courseID)

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))

	result, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userObjID},
		bson.M{"$pull": bson.M{"cart": courseObjID}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from cart"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course not in cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func Checkout(c *gin.Context) {
	var input struct {
		IDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))

	var courseIDs []primitive.ObjectID
	for _, id := range input.IDs {
		objID, _ := primitive.ObjectIDFromHex(id)
		courseIDs = append(courseIDs, objID)
	}

	enrolledCourses := make([]models.EnrolledCourse, len(courseIDs))
	for i, courseID := range courseIDs {
		enrolledCourses[i] = models.EnrolledCourse{
			Course:     courseID,
			EnrolledOn: time.Now(),
		}
	}

	// Update courses with enrollments
	_, err := database.DB.Collection("courses").UpdateMany(
		context.Background(),
		bson.M{
			"_id": bson.M{"$in": courseIDs},
			"pricing": "Free",
		},
		bson.M{
			"$addToSet": bson.M{
				"meta.enrollments": bson.M{
					"id": userObjID,
				},
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update courses"})
		return
	}

	// Update user enrollments and clear cart/wishlist
	result, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userObjID},
		bson.M{
			"$addToSet": bson.M{"enrolledCourses": bson.M{"$each": enrolledCourses}},
			"$pull": bson.M{
				"cart": bson.M{"$in": courseIDs},
				"wishlist": bson.M{"$in": courseIDs},
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, result)
}