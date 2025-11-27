package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// FavoritesQuery represents query parameters for favorites endpoints
type FavoritesQuery struct {
	PaginationQuery
	Type string `form:"type" binding:"required,oneof=exercise workout"`
}

// FavoritesTypeQuery for endpoints that only need type
type FavoritesTypeQuery struct {
	Type string `form:"type" binding:"required,oneof=exercise workout"`
}

// GetFavorites returns user's favorites filtered by type
// GET /api/v1/user/favorites?type=exercise|workout
func GetFavorites(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var params FavoritesQuery
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	SetDefaultPagination(&params.PaginationQuery)
	offset := (params.Page - 1) * params.Limit

	switch params.Type {
	case "exercise":
		getFavoriteExercises(c, userID, offset, params.Limit, params.Page)
	case "workout":
		getFavoriteWorkouts(c, userID, offset, params.Limit, params.Page)
	}
}

func getFavoriteExercises(c *gin.Context, userID uuid.UUID, offset, limit, page int) {
	var total int64
	if err := database.DB.Model(&models.UserFavoriteExercise{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count favorite exercises.")
		return
	}

	var favorites []models.UserFavoriteExercise
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Exercise").
		Preload("Exercise.MuscleGroups.MuscleGroup").
		Preload("Exercise.Equipment.Equipment").
		Preload("Exercise.ExerciseTypes.ExerciseType").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&favorites).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch favorite exercises.")
		return
	}

	responses := make([]models.FavoriteResponse, len(favorites))
	for i, fav := range favorites {
		responses[i] = fav.ToGenericResponse()
	}

	utils.PaginatedResponse(c, "Favorite exercises retrieved successfully.", responses, page, limit, int(total))
}

func getFavoriteWorkouts(c *gin.Context, userID uuid.UUID, offset, limit, page int) {
	var total int64
	if err := database.DB.Model(&models.UserFavoriteWorkout{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count favorite workouts.")
		return
	}

	var favorites []models.UserFavoriteWorkout
	if err := database.DB.Where("user_id = ?", userID).
		Preload("Workout").
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&favorites).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch favorite workouts.")
		return
	}

	responses := make([]models.FavoriteResponse, len(favorites))
	for i, fav := range favorites {
		responses[i] = fav.ToGenericResponse()
	}

	utils.PaginatedResponse(c, "Favorite workouts retrieved successfully.", responses, page, limit, int(total))
}

// AddFavorite adds an item to user's favorites
// POST /api/v1/user/favorites?type=exercise|workout
func AddFavorite(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var typeQuery FavoritesTypeQuery
	if err := c.ShouldBindQuery(&typeQuery); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var req models.AddFavoriteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	switch typeQuery.Type {
	case "exercise":
		addFavoriteExercise(c, userID, req.ItemID)
	case "workout":
		addFavoriteWorkout(c, userID, req.ItemID)
	}
}

func addFavoriteExercise(c *gin.Context, userID, exerciseID uuid.UUID) {
	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.First(&exercise, "id = ?", exerciseID).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	// Check if already favorited
	var existing models.UserFavoriteExercise
	if err := database.DB.Where("user_id = ? AND exercise_id = ?", userID, exerciseID).
		First(&existing).Error; err == nil {
		utils.ConflictResponse(c, "Exercise is already in favorites.")
		return
	}

	favorite := models.UserFavoriteExercise{
		UserID:     userID,
		ExerciseID: exerciseID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add exercise to favorites.")
		return
	}

	// Load exercise for response
	database.DB.Preload("Exercise").
		Preload("Exercise.MuscleGroups.MuscleGroup").
		Preload("Exercise.Equipment.Equipment").
		Preload("Exercise.ExerciseTypes.ExerciseType").
		First(&favorite, favorite.ID)

	utils.CreatedResponse(c, "Exercise added to favorites.", favorite.ToGenericResponse())
}

func addFavoriteWorkout(c *gin.Context, userID, workoutID uuid.UUID) {
	// Check if workout exists
	var workout models.Workout
	if err := database.DB.First(&workout, "id = ?", workoutID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Check visibility
	if !canUserViewWorkout(userID, &workout) {
		utils.NotFoundResponse(c, "Workout not found.")
		return
	}

	// Check if already favorited
	var existing models.UserFavoriteWorkout
	if err := database.DB.Where("user_id = ? AND workout_id = ?", userID, workoutID).
		First(&existing).Error; err == nil {
		utils.ConflictResponse(c, "Workout is already in favorites.")
		return
	}

	favorite := models.UserFavoriteWorkout{
		UserID:    userID,
		WorkoutID: workoutID,
	}

	if err := database.DB.Create(&favorite).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to add workout to favorites.")
		return
	}

	// Load workout for response
	database.DB.Preload("Workout").First(&favorite, favorite.ID)

	utils.CreatedResponse(c, "Workout added to favorites.", favorite.ToGenericResponse())
}

// RemoveFavorite removes an item from user's favorites
// DELETE /api/v1/user/favorites/:id?type=exercise|workout
func RemoveFavorite(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var typeQuery FavoritesTypeQuery
	if err := c.ShouldBindQuery(&typeQuery); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	itemID, ok := utils.ParseUUID(c, c.Param("id"), "item")
	if !ok {
		return
	}

	switch typeQuery.Type {
	case "exercise":
		removeFavoriteExercise(c, userID, itemID)
	case "workout":
		removeFavoriteWorkout(c, userID, itemID)
	}
}

func removeFavoriteExercise(c *gin.Context, userID, exerciseID uuid.UUID) {
	result := database.DB.Where("user_id = ? AND exercise_id = ?", userID, exerciseID).
		Delete(&models.UserFavoriteExercise{})

	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove exercise from favorites.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Favorite not found.")
		return
	}

	utils.DeletedResponse(c, "Exercise removed from favorites.")
}

func removeFavoriteWorkout(c *gin.Context, userID, workoutID uuid.UUID) {
	result := database.DB.Where("user_id = ? AND workout_id = ?", userID, workoutID).
		Delete(&models.UserFavoriteWorkout{})

	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove workout from favorites.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Favorite not found.")
		return
	}

	utils.DeletedResponse(c, "Workout removed from favorites.")
}

// CheckFavorite checks if an item is favorited by the user
// GET /api/v1/user/favorites/:id/check?type=exercise|workout
func CheckFavorite(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var typeQuery FavoritesTypeQuery
	if err := c.ShouldBindQuery(&typeQuery); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	itemID, ok := utils.ParseUUID(c, c.Param("id"), "item")
	if !ok {
		return
	}

	var isFavorited bool

	switch typeQuery.Type {
	case "exercise":
		var count int64
		database.DB.Model(&models.UserFavoriteExercise{}).
			Where("user_id = ? AND exercise_id = ?", userID, itemID).
			Count(&count)
		isFavorited = count > 0
	case "workout":
		var count int64
		database.DB.Model(&models.UserFavoriteWorkout{}).
			Where("user_id = ? AND workout_id = ?", userID, itemID).
			Count(&count)
		isFavorited = count > 0
	}

	utils.SuccessResponse(c, "Favorite status retrieved.", gin.H{
		"is_favorited": isFavorited,
	})
}

// canUserViewWorkout checks if a user can view a workout based on visibility rules
func canUserViewWorkout(userID uuid.UUID, workout *models.Workout) bool {
	// User owns the workout
	if workout.UserID == userID {
		return true
	}

	// Public workout
	if workout.Visibility == "public" {
		return true
	}

	// Friends-only workout - check friendship
	if workout.Visibility == "friends" {
		var count int64
		database.DB.Model(&models.Friendship{}).
			Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?",
				userID, workout.UserID, workout.UserID, userID, "accepted").
			Count(&count)
		return count > 0
	}

	return false
}
