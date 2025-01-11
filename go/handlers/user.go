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
