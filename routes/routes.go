package routes

import (
	"lamari-fit-api/controllers"
	"lamari-fit-api/middleware"
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

		// Public invitation verification (no auth required)
		api.GET("/invitations/verify/:token", controllers.VerifyInvitationToken)

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

			// Exercise Types
			exerciseTypes := protected.Group("/exercise-types")
			{
				exerciseTypes.GET("/", controllers.GetExerciseTypes)
				exerciseTypes.GET("/:id", controllers.GetExerciseType)
				exerciseTypes.POST("/", controllers.CreateExerciseType)
				exerciseTypes.PUT("/:id", controllers.UpdateExerciseType)
				exerciseTypes.DELETE("/:id", controllers.DeleteExerciseType)
			}

			// Exercises
			exercises := protected.Group("/exercises")
			{
				exercises.POST("/", controllers.CreateExercise)
				exercises.GET("/", controllers.GetExercises)
				exercises.GET("/by-slug/:slug", controllers.GetExerciseBySlug)
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

				// Exercise-Type relationships
				exercises.POST("/:id/types", controllers.AssignExerciseType)
				exercises.GET("/:id/types", controllers.GetExerciseTypesByExercise)
				exercises.DELETE("/:id/types/:type_id", controllers.RemoveExerciseType)
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

			// User Settings
			userSettings := protected.Group("/user/settings")
			{
				userSettings.GET("", controllers.GetUserSettings)
				userSettings.PUT("", controllers.UpdateUserSettings)
			}

			// User Fitness Settings
			userFitness := protected.Group("/user/fitness")
			{
				userFitness.PUT("/level", controllers.UpdateUserFitnessLevel)
			}

			// User Fitness Profile
			userFitnessProfile := protected.Group("/user/fitness-profile")
			{
				userFitnessProfile.POST("", controllers.CreateUserFitnessProfile)
				userFitnessProfile.GET("", controllers.GetUserFitnessProfile)
				userFitnessProfile.PUT("", controllers.UpdateUserFitnessProfile)
				userFitnessProfile.DELETE("", controllers.DeleteUserFitnessProfile)
				userFitnessProfile.POST("/log-weight", controllers.LogWeight)
			}

			// Weight Logs
			weightLogs := protected.Group("/user/weight-logs")
			{
				weightLogs.POST("", controllers.CreateWeightLog)
				weightLogs.GET("", controllers.GetWeightLogs)
				weightLogs.GET("/stats", controllers.GetWeightStats)
				weightLogs.GET("/:id", controllers.GetWeightLog)
				weightLogs.PUT("/:id", controllers.UpdateWeightLog)
				weightLogs.DELETE("/:id", controllers.DeleteWeightLog)
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

			// User Favorites
			userFavorites := protected.Group("/user/favorites")
			{
				userFavorites.GET("", controllers.GetFavorites)
				userFavorites.POST("", controllers.AddFavorite)
				userFavorites.DELETE("/:id", controllers.RemoveFavorite)
				userFavorites.GET("/:id/check", controllers.CheckFavorite)
			}

			// Workout Plans
			workoutPlans := protected.Group("/workout-plans")
			{
				workoutPlans.POST("/", controllers.CreateWorkoutPlan)
				workoutPlans.GET("/", controllers.GetWorkoutPlans)
				workoutPlans.GET("/:id", controllers.GetWorkoutPlan)
				workoutPlans.PUT("/:id", controllers.UpdateWorkoutPlan)
				workoutPlans.DELETE("/:id", controllers.DeleteWorkoutPlan)

				// Workout Plan Items
				workoutPlans.POST("/:id/workouts", controllers.AddWorkoutToPlan)
				workoutPlans.GET("/:id/workouts", controllers.GetPlanWorkouts)
				workoutPlans.DELETE("/:id/workouts/:item_id", controllers.RemoveWorkoutFromPlan)
				workoutPlans.PUT("/:id/workouts/:item_id", controllers.UpdatePlanItem)
			}

			// Workouts
			workouts := protected.Group("/workouts")
			{
				workouts.POST("/", controllers.CreateWorkout)
				workouts.GET("/", controllers.GetUserWorkouts)
				workouts.GET("/:id", controllers.GetWorkout)
				workouts.PUT("/:id", controllers.UpdateWorkout)
				workouts.DELETE("/:id", controllers.DeleteWorkout)
				workouts.POST("/:id/duplicate", controllers.DuplicateWorkout)

				// Workout Prescriptions
				workouts.POST("/:id/prescriptions", controllers.CreatePrescriptionGroup)
				workouts.GET("/:id/prescriptions", controllers.GetWorkoutPrescriptions)
				workouts.PUT("/:id/prescriptions/reorder", controllers.ReorderPrescriptionGroups)
				workouts.PUT("/:id/prescriptions/:group_id", controllers.UpdatePrescriptionGroup)
				workouts.DELETE("/:id/prescriptions/:group_id", controllers.DeletePrescriptionGroup)
				workouts.POST("/:id/prescriptions/:group_id/exercises", controllers.AddExerciseToPrescriptionGroup)
			}

			// Plan Enrollments
			enrollments := protected.Group("/enrollments")
			{
				enrollments.POST("/", controllers.EnrollInPlan)
				enrollments.GET("/", controllers.GetUserEnrollments)
				enrollments.GET("/:id", controllers.GetEnrollment)
				enrollments.PUT("/:id", controllers.UpdateEnrollment)
				enrollments.DELETE("/:id", controllers.CancelEnrollment)
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

			// Trainers
			trainers := protected.Group("/trainers")
			{
				// Profile management (authenticated user's own profile)
				trainers.POST("/profile", controllers.CreateTrainerProfile)
				trainers.GET("/profile", controllers.GetTrainerProfile)
				trainers.PUT("/profile", controllers.UpdateTrainerProfile)
				trainers.DELETE("/profile", controllers.DeleteTrainerProfile)

				// Client management (trainer side)
				trainers.POST("/clients", controllers.InviteClient)
				trainers.GET("/clients", controllers.GetTrainerClients)
				trainers.DELETE("/clients/:id", controllers.RemoveClient)

				// Email invitations (trainer side)
				trainers.POST("/email-invitations", controllers.CreateEmailInvitation)
				trainers.GET("/email-invitations", controllers.GetEmailInvitations)
				trainers.DELETE("/email-invitations/:id", controllers.CancelEmailInvitation)
				trainers.POST("/email-invitations/:id/resend", controllers.ResendEmailInvitation)

				// Public trainer endpoints
				trainers.GET("/", controllers.ListTrainers)
				trainers.GET("/:id", controllers.GetTrainerPublicProfile)
			}

			// User's trainer relationships (client side)
			me := protected.Group("/me")
			{
				me.GET("/trainers", controllers.GetMyTrainers)
				me.GET("/trainer-invitations", controllers.GetMyTrainerInvitations)
				me.PUT("/trainer-invitations/:id", controllers.RespondToInvitation)
			}

			// Specialties
			specialties := protected.Group("/specialties")
			{
				specialties.GET("/", controllers.ListSpecialties)
			}

			// RPE Scales
			rpe := protected.Group("/rpe")
			{
				rpe.GET("/scales", controllers.ListRPEScales)
				rpe.GET("/scales/global", controllers.GetGlobalRPEScale)
				rpe.POST("/scales", controllers.CreateRPEScale)
				rpe.GET("/scales/:id", controllers.GetRPEScale)
				rpe.PUT("/scales/:id", controllers.UpdateRPEScale)
				rpe.DELETE("/scales/:id", controllers.DeleteRPEScale)
				rpe.POST("/scales/:id/values", controllers.AddRPEScaleValue)
			}

			// Workout Sessions (Logging)
			workoutSessions := protected.Group("/workout-sessions")
			{
				workoutSessions.POST("", controllers.CreateWorkoutSession)
				workoutSessions.GET("", controllers.GetWorkoutSessions)
				workoutSessions.GET("/:id", controllers.GetWorkoutSession)
				workoutSessions.PUT("/:id", controllers.UpdateWorkoutSession)
				workoutSessions.PUT("/:id/end", controllers.EndWorkoutSession)
				workoutSessions.DELETE("/:id", controllers.DeleteWorkoutSession)
			}

			// Session Blocks
			sessionBlocks := protected.Group("/session-blocks")
			{
				sessionBlocks.GET("/:id", controllers.GetSessionBlock)
				sessionBlocks.PUT("/:id/complete", controllers.CompleteSessionBlock)
				sessionBlocks.PUT("/:id/skip", controllers.SkipSessionBlock)
				sessionBlocks.PUT("/:id/rpe", controllers.UpdateSessionBlockRPE)
			}

			// Session Exercises
			sessionExercises := protected.Group("/session-exercises")
			{
				sessionExercises.GET("/:id", controllers.GetSessionExercise)
				sessionExercises.PUT("/:id/complete", controllers.CompleteSessionExercise)
				sessionExercises.PUT("/:id/skip", controllers.SkipSessionExercise)
				sessionExercises.PUT("/:id/notes", controllers.UpdateSessionExerciseNotes)
				sessionExercises.POST("/:id/sets", controllers.AddSetToExercise)
			}

			// Session Sets
			sessionSets := protected.Group("/session-sets")
			{
				sessionSets.GET("/:id", controllers.GetSessionSet)
				sessionSets.PUT("/:id", controllers.UpdateSessionSet)
				sessionSets.PUT("/:id/complete", controllers.CompleteSessionSet)
				sessionSets.DELETE("/:id", controllers.DeleteSessionSet)
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

			protected.GET("/nutrition", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Your nutrition data",
					"meals":   []string{"Breakfast", "Lunch", "Dinner"},
				})
			})
		}
	}
}
