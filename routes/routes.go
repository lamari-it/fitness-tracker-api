package routes

import (
	"fit-flow-api/controllers"
	"fit-flow-api/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Language")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		
		c.Next()
	})
	
	// Add i18n middleware
	r.Use(middleware.I18nMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := r.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.GET("/google", controllers.GoogleLogin)
			auth.GET("/google/callback", controllers.GoogleCallback)
			auth.POST("/apple", controllers.AppleLogin)
			
			auth.Use(middleware.AuthMiddleware())
			auth.GET("/profile", controllers.GetProfile)
		}

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/dashboard", func(c *gin.Context) {
				userID := c.GetString("user_id")
				email := c.GetString("email")
				
				c.JSON(http.StatusOK, gin.H{
					"message": "Welcome to your dashboard",
					"user_id": userID,
					"email":   email,
				})
			})
			
			// Muscle Groups
			muscleGroups := protected.Group("/muscle-groups")
			{
				muscleGroups.POST("/", controllers.CreateMuscleGroup)
				muscleGroups.GET("/", controllers.GetMuscleGroups)
				muscleGroups.GET("/:id", controllers.GetMuscleGroup)
				muscleGroups.PUT("/:id", controllers.UpdateMuscleGroup)
				muscleGroups.DELETE("/:id", controllers.DeleteMuscleGroup)
			}

			// Exercises
			exercises := protected.Group("/exercises")
			{
				exercises.POST("/", controllers.CreateExercise)
				exercises.GET("/", controllers.GetExercises)
				exercises.GET("/:id", controllers.GetExercise)
				exercises.PUT("/:id", controllers.UpdateExercise)
				exercises.DELETE("/:id", controllers.DeleteExercise)
				
				// Exercise-MuscleGroup relationships
				exercises.POST("/:id/muscle-groups", controllers.AssignMuscleGroupToExercise)
				exercises.GET("/:id/muscle-groups", controllers.GetExerciseMuscleGroups)
				exercises.DELETE("/:id/muscle-groups/:muscle_group_id", controllers.RemoveMuscleGroupFromExercise)
				
				// Exercise-Equipment relationships
				exercises.POST("/:id/equipment", controllers.AssignEquipmentToExercise)
				exercises.GET("/:id/equipment", controllers.GetExerciseEquipment)
				exercises.DELETE("/:id/equipment/:equipment_id", controllers.RemoveEquipmentFromExercise)
			}

			// Equipment
			equipment := protected.Group("/equipment")
			{
				equipment.POST("/", controllers.CreateEquipment)
				equipment.GET("/", controllers.GetAllEquipment)
				equipment.GET("/:id", controllers.GetEquipmentByID)
				equipment.PUT("/:id", controllers.UpdateEquipment)
				equipment.DELETE("/:id", controllers.DeleteEquipment)
			}

			// Fitness Levels
			fitnessLevels := protected.Group("/fitness-levels")
			{
				fitnessLevels.GET("/", controllers.GetAllFitnessLevels)
				fitnessLevels.GET("/:id", controllers.GetFitnessLevel)
				fitnessLevels.POST("/", controllers.CreateFitnessLevel)
				fitnessLevels.PUT("/:id", controllers.UpdateFitnessLevel)
				fitnessLevels.DELETE("/:id", controllers.DeleteFitnessLevel)
			}

			// Fitness Goals
			fitnessGoals := protected.Group("/fitness-goals")
			{
				fitnessGoals.GET("/", controllers.GetAllFitnessGoals)
				fitnessGoals.GET("/:id", controllers.GetFitnessGoal)
				fitnessGoals.POST("/", controllers.CreateFitnessGoal)
				fitnessGoals.PUT("/:id", controllers.UpdateFitnessGoal)
				fitnessGoals.DELETE("/:id", controllers.DeleteFitnessGoal)
			}

			// User Fitness Settings
			userFitness := protected.Group("/user/fitness")
			{
				userFitness.GET("/goals", controllers.GetUserFitnessGoals)
				userFitness.PUT("/goals", controllers.SetUserFitnessGoals)
				userFitness.PUT("/level", controllers.UpdateUserFitnessLevel)
			}

			// User Equipment
			userEquipment := protected.Group("/user/equipment")
			{
				userEquipment.GET("/", controllers.GetUserEquipment)
				userEquipment.POST("/", controllers.AddUserEquipment)
				userEquipment.POST("/bulk", controllers.BulkAddUserEquipment)
				userEquipment.PUT("/:id", controllers.UpdateUserEquipment)
				userEquipment.DELETE("/:id", controllers.RemoveUserEquipment)
				userEquipment.GET("/location/:location", controllers.GetUserEquipmentByLocation)
			}

			// Workout Plans
			workoutPlans := protected.Group("/workout-plans")
			{
				workoutPlans.POST("/", controllers.CreateWorkoutPlan)
				workoutPlans.GET("/", controllers.GetWorkoutPlans)
				workoutPlans.GET("/:id", controllers.GetWorkoutPlan)
				workoutPlans.PUT("/:id", controllers.UpdateWorkoutPlan)
				workoutPlans.DELETE("/:id", controllers.DeleteWorkoutPlan)
			}

			// Friends
			friends := protected.Group("/friends")
			{
				friends.POST("/request", controllers.SendFriendRequest)
				friends.GET("/requests", controllers.GetFriendRequests)
				friends.PUT("/requests/:id/:action", controllers.RespondToFriendRequest)
				friends.GET("/", controllers.GetFriends)
				friends.DELETE("/:id", controllers.RemoveFriend)
			}
			
			// Translations (Admin only)
			translations := protected.Group("/translations")
			translations.Use(middleware.AdminMiddleware())
			{
				translations.POST("/", controllers.CreateTranslation)
				translations.GET("/", controllers.GetTranslations)
				translations.GET("/:id", controllers.GetTranslation)
				translations.PUT("/:id", controllers.UpdateTranslation)
				translations.DELETE("/:id", controllers.DeleteTranslation)
				translations.GET("/resource/:resource_type/:resource_id", controllers.GetResourceTranslations)
				translations.POST("/upsert", controllers.CreateOrUpdateTranslation)
			}

			// Legacy endpoints for backwards compatibility
			protected.GET("/workouts", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Your workouts",
					"workouts": []string{"Morning Run", "Evening Yoga", "Strength Training"},
				})
			})
			
			protected.GET("/nutrition", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Your nutrition data",
					"meals": []string{"Breakfast", "Lunch", "Dinner"},
				})
			})
		}
	}
}