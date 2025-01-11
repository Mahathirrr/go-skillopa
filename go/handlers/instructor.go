package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"learnlit/database"
	"learnlit/models"
)

func MakeInstructor(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	// Get user details
	var user models.User
	err := database.DB.Collection("users").FindOne(context.Background(), 
		bson.M{"_id": objID}).Decode(&user)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already instructor
	for _, role := range user.Role {
		if role == "Instructor" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Already an Instructor"})
			return
		}
	}

	// Create instructor profile
	instructor := models.Instructor{
		Name:  user.Name,
		Email: user.Email,
	}

	result, err := database.DB.Collection("instructors").InsertOne(context.Background(), instructor)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create instructor profile"})
		return
	}

	// Update user role
	_, err = database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{
			"$set": bson.M{
				"instructorProfile": result.InsertedID,
				"role":             append(user.Role, "Instructor"),
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	c.JSON(http.StatusCreated, instructor)
}

func GetInstructorProfile(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var instructor models.Instructor
	err = database.DB.Collection("instructors").FindOne(context.Background(), 
		bson.M{"_id": objID}).Decode(&instructor)
	
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instructor not found"})
		return
	}

	c.JSON(http.StatusOK, instructor)
}

func UpdateInstructorProfile(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid instructor ID"})
		return
	}

	var input models.Instructor
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := database.DB.Collection("instructors").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": input},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update instructor"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instructor not found"})
		return
	}

	c.JSON(http.StatusOK, input)
}