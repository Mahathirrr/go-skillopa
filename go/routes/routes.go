package routes

import (
	"learnlit/handlers"
	"learnlit/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	api := r.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	{
		auth.POST("/register", handlers.Register)
		auth.POST("/login", handlers.Login)
		auth.POST("/google-login", handlers.GoogleLogin)
		auth.POST("/logout", handlers.Logout)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware())
	{
		// User routes
		protected.GET("/user/current-user", handlers.CurrentUser)
		protected.GET("/user/cart", handlers.GetCart)
		protected.POST("/user/cart", handlers.AddToCart)
		protected.DELETE("/user/cart/:id", handlers.RemoveFromCart)
		protected.GET("/user/wishlist", handlers.GetWishlist)
		protected.POST("/user/wishlist", handlers.AddToWishlist)
		protected.DELETE("/user/wishlist/:id", handlers.RemoveFromWishlist)
		protected.POST("/checkout", handlers.Checkout)
		protected.GET("/user/enrolled-courses", handlers.GetEnrolledCourses)
		protected.PUT("/user/profile", handlers.UpdateProfile)

		// Course routes
		protected.POST("/create-course", handlers.CreateCourse)
		protected.PUT("/course/:id", handlers.UpdateCourse)
		protected.GET("/me/taught-courses", handlers.GetTaughtCourses)
		protected.GET("/me/posted-courses", handlers.GetPostedCourses)
	}

	// Public course routes
	courses := api.Group("/courses")
	{
		courses.GET("", handlers.GetCategoryCourses)
		courses.GET("/search", handlers.SearchCourses)
		courses.GET("/all-courses", handlers.GetAllPublishedCourses)
		courses.GET("/course-categories", handlers.GetCourseCategories)
		courses.POST("/get-course", handlers.GetCourse)
	}

	// Payment routes
	payments := api.Group("/payment")
	payments.Use(middleware.AuthMiddleware())
	{
		payments.POST("/create", handlers.CreatePayment)
		payments.GET("/status/:orderId", handlers.GetPaymentStatus)
		payments.POST("/notification", handlers.HandlePaymentNotification)
	}
}
