package handlers

import (
	"context"
	"learnlit/database"
	"learnlit/models"
	"learnlit/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllPublishedCourses(c *gin.Context) {
	courses := []models.Course{}

	opts := options.Find().SetProjection(bson.M{
		"meta":       0,
		"curriculum": 0,
	})

	cursor, err := database.DB.Collection("courses").Find(context.Background(),
		bson.M{"published": true}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &courses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func CreateCourse(c *gin.Context) {
	var input models.Course
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	// Create slug
	slug := utils.CreateSlug(input.Title)

	// Check if title exists
	var existingCourse models.Course
	err := database.DB.Collection("courses").FindOne(context.Background(),
		bson.M{"slug": slug}).Decode(&existingCourse)

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Title is taken"})
		return
	}

	input.Slug = slug
	input.PostedBy = objID

	result, err := database.DB.Collection("courses").InsertOne(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}

func GetCourse(c *gin.Context) {
	var input struct {
		ID     string `json:"id"`
		Slug   string `json:"slug"`
		UserID string `json:"userId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var filter primitive.M
	if input.ID != "" {
		objID, err := primitive.ObjectIDFromHex(input.ID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		filter = bson.M{"_id": objID}
	} else if input.Slug != "" {
		filter = bson.M{"slug": input.Slug}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing course id/slug"})
		return
	}

	var course models.Course
	err := database.DB.Collection("courses").FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	if input.UserID != "" {
		userObjID, _ := primitive.ObjectIDFromHex(input.UserID)
		for _, enrollment := range course.Meta.Enrollments {
			if enrollment.ID == userObjID {
				course.Meta.Enrollments = nil
				c.JSON(http.StatusOK, gin.H{
					"course":         course,
					"isUserEnrolled": true,
				})
				return
			}
		}
	}

	course.Meta.Enrollments = nil
	c.JSON(http.StatusOK, gin.H{
		"course":         course,
		"isUserEnrolled": false,
	})
}

