package handlers

import (
	"context"
	"learnlit/database"
	"learnlit/models"
	"learnlit/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CurrentUser(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var user models.User
	err := database.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	user.Password = "" // Don't send password
	c.JSON(http.StatusOK, user)
}

func GetCart(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var user models.User
	err := database.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Cart)
}

func GetWishlist(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var user models.User
	err := database.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.Wishlist)
}

func AddToWishlist(c *gin.Context) {
	var input struct {
		CourseID string `json:"courseId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))
	courseID, _ := primitive.ObjectIDFromHex(input.CourseID)

	_, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$addToSet": bson.M{"wishlist": courseID}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to wishlist"})
}

func RemoveFromWishlist(c *gin.Context) {
	courseID := c.Param("id")
	courseObjID, _ := primitive.ObjectIDFromHex(courseID)

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))

	_, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userObjID},
		bson.M{"$pull": bson.M{"wishlist": courseObjID}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove from wishlist"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from wishlist"})
}

func GetEnrolledCourses(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var user models.User
	err := database.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": objID}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user.EnrolledCourses)
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userId")
	objID, _ := primitive.ObjectIDFromHex(userID.(string))

	var input struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{}
	if input.Name != "" {
		update["name"] = input.Name
	}
	if input.Avatar != "" {
		// Handle avatar upload using cloudinary
		avatarURL, err := utils.UploadToCloudinary([]byte(input.Avatar))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload avatar"})
			return
		}
		update["avatar"] = avatarURL
	}

	_, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func AddToCart(c *gin.Context) {
	var input struct {
		CourseID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))
	courseObjID, _ := primitive.ObjectIDFromHex(input.CourseID)

	result, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userObjID},
		bson.M{"$addToSet": bson.M{"cart": courseObjID}},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to cart"})
		return
	}

	if result.ModifiedCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Course already in cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Added to cart"})
}

func RemoveFromCart(c *gin.Context) {
	courseID := c.Param("id")
	courseObjID, err := primitive.ObjectIDFromHex(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

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
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found in cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Removed from cart"})
}

func Checkout(c *gin.Context) {
	var input struct {
		CourseIDs []string `json:"ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, _ := c.Get("userId")
	userObjID, _ := primitive.ObjectIDFromHex(userID.(string))

	var courseObjIDs []primitive.ObjectID
	for _, id := range input.CourseIDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
			return
		}
		courseObjIDs = append(courseObjIDs, objID)
	}

	// Add courses to enrolled courses
	enrolledCourses := make([]models.EnrolledCourse, len(courseObjIDs))
	for i, courseID := range courseObjIDs {
		enrolledCourses[i] = models.EnrolledCourse{
			Course:     courseID,
			EnrolledOn: time.Now(),
		}
	}

	_, err := database.DB.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"_id": userObjID},
		bson.M{
			"$addToSet": bson.M{"enrolledCourses": bson.M{"$each": enrolledCourses}},
			"$pull": bson.M{
				"cart":     bson.M{"$in": courseObjIDs},
				"wishlist": bson.M{"$in": courseObjIDs},
			},
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process checkout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Checkout successful"})
}
