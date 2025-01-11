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

func GetAllPublishedCourses(c *gin.Context) {
	var courses []models.Course
	cursor, err := database.DB.Collection("courses").Find(context.Background(), bson.M{"published": true})
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

func GetCourse(c *gin.Context) {
	var input struct {
		ID   string `json:"id"`
		Slug string `json:"slug"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var filter bson.M
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID or slug is required"})
		return
	}

	var course models.Course
	err := database.DB.Collection("courses").FindOne(context.Background(), filter).Decode(&course)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))
	course.PostedBy = objID

	result, err := database.DB.Collection("courses").InsertOne(context.Background(), course)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		course.ID = oid
	}

	c.JSON(http.StatusCreated, course)
}

func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	var updates models.Course
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))

	result, err := database.DB.Collection("courses").UpdateOne(
		context.Background(),
		bson.M{"_id": objID, "postedBy": userObjID},
		bson.M{"$set": updates},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update course"})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}
