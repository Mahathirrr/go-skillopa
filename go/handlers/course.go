package handlers

import (
	"context"
	"learnlit/data"
	"learnlit/database"
	"learnlit/models"
	"learnlit/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCourseCategories(c *gin.Context) {
	c.JSON(http.StatusOK, data.Categories)
}

func GetCategoryCourses(c *gin.Context) {
	category := c.Query("category")
	subCategory := c.Query("subCategory")

	filter := bson.M{"published": true}
	if category != "" {
		filter["category"] = utils.SnakeCaseToTitle(category)
	}
	if subCategory != "" {
		filter["subCategory"] = utils.SnakeCaseToTitle(subCategory)
	}

	var courses []models.Course
	cursor, err := database.DB.Collection("courses").Find(context.Background(), filter)
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

func SearchCourses(c *gin.Context) {
	query := c.Query("query")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}

	filter := bson.M{
		"published": true,
		"$text":     bson.M{"$search": query},
	}

	var courses []models.Course
	cursor, err := database.DB.Collection("courses").Find(context.Background(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search courses"})
		return
	}
	defer cursor.Close(context.Background())

	if err = cursor.All(context.Background(), &courses); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}

func GetTaughtCourses(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var instructor models.Instructor
	err := database.DB.Collection("instructors").FindOne(context.Background(),
		bson.M{"_id": objID}).Decode(&instructor)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Instructor not found"})
		return
	}

	var courses []models.Course
	cursor, err := database.DB.Collection("courses").Find(context.Background(),
		bson.M{"_id": bson.M{"$in": instructor.Courses}})
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

func GetPostedCourses(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var courses []models.Course
	cursor, err := database.DB.Collection("courses").Find(context.Background(),
		bson.M{"postedBy": objID})
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
