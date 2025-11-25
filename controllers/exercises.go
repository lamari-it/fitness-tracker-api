package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// generateSlug creates a URL-friendly slug from a name
func generateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)
	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	// Remove apostrophes
	slug = strings.ReplaceAll(slug, "'", "")
	// Replace multiple hyphens with single hyphen
	slug = strings.ReplaceAll(slug, "--", "-")
	return slug
}

type CreateExerciseRequest struct {
	Name         string                  `json:"name" binding:"required"`
	Description  string                  `json:"description"`
	Equipment    string                  `json:"equipment"`
	IsBodyweight bool                    `json:"is_bodyweight"`
	Instructions string                  `json:"instructions"`
	VideoURL     string                  `json:"video_url"`
	MuscleGroups []MuscleGroupAssignment `json:"muscle_groups,omitempty"`
}

type MuscleGroupAssignment struct {
	MuscleGroupID uuid.UUID `json:"muscle_group_id" binding:"required"`
	Primary       bool      `json:"primary"`
	Intensity     string    `json:"intensity"`
}

func CreateExercise(c *gin.Context) {
	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Start transaction
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	exercise := models.Exercise{
		Slug:         generateSlug(req.Name),
		Name:         req.Name,
		Description:  req.Description,
		IsBodyweight: req.IsBodyweight,
		Instructions: req.Instructions,
		VideoURL:     req.VideoURL,
	}

	if err := tx.Create(&exercise).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key") {
			utils.ConflictResponse(c, "Exercise with this name already exists.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to create exercise.")
		return
	}

	// Assign muscle groups if provided
	if len(req.MuscleGroups) > 0 {
		primaryCount := 0
		for _, mgAssign := range req.MuscleGroups {
			if mgAssign.Primary {
				primaryCount++
			}
		}

		// Ensure only one primary muscle group
		if primaryCount > 1 {
			tx.Rollback()
			utils.BadRequestResponse(c, "Only one muscle group can be set as primary.", nil)
			return
		}

		for _, mgAssign := range req.MuscleGroups {
			// Verify muscle group exists
			var muscleGroup models.MuscleGroup
			if err := tx.Where("id = ?", mgAssign.MuscleGroupID).First(&muscleGroup).Error; err != nil {
				tx.Rollback()
				utils.BadRequestResponse(c, "Invalid muscle group ID: "+mgAssign.MuscleGroupID.String(), nil)
				return
			}

			assignment := models.ExerciseMuscleGroup{
				ExerciseID:    exercise.ID,
				MuscleGroupID: mgAssign.MuscleGroupID,
				Primary:       mgAssign.Primary,
				Intensity:     mgAssign.Intensity,
			}

			if assignment.Intensity == "" {
				assignment.Intensity = "moderate"
			}

			if err := assignment.Validate(); err != nil {
				tx.Rollback()
				utils.BadRequestResponse(c, "Invalid muscle group assignment.", err.Error())
				return
			}

			if err := tx.Create(&assignment).Error; err != nil {
				tx.Rollback()
				utils.InternalServerErrorResponse(c, "Failed to assign muscle groups.")
				return
			}
		}
	}

	tx.Commit()

	// Load the exercise with muscle groups for response
	database.DB.Where("id = ?", exercise.ID).
		Preload("MuscleGroups.MuscleGroup").
		First(&exercise)

	utils.CreatedResponse(c, "Exercise created successfully.", exercise)
}

func GetExercises(c *gin.Context) {
	var params ExerciseQuery
	if err := c.ShouldBindQuery(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&params.PaginationQuery)

	offset := (params.Page - 1) * params.Limit

	query := database.DB.Model(&models.Exercise{}).
		Preload("MuscleGroups.MuscleGroup")

	if params.Search != "" {
		query = query.Where("name ILIKE ?", "%"+params.Search+"%")
	}

	if params.MuscleGroupID != "" {
		// Filter by muscle group through the many-to-many relationship
		if params.PrimaryOnly == "true" {
			query = query.Joins("JOIN exercise_muscle_groups emg ON exercises.id = emg.exercise_id").
				Where("emg.muscle_group_id = ? AND emg.primary = true", params.MuscleGroupID)
		} else {
			query = query.Joins("JOIN exercise_muscle_groups emg ON exercises.id = emg.exercise_id").
				Where("emg.muscle_group_id = ?", params.MuscleGroupID)
		}
	}

	if params.Equipment != "" {
		query = query.Where("equipment = ?", params.Equipment)
	}

	if params.Bodyweight != "" {
		isBodyweight := params.Bodyweight == "true"
		query = query.Where("is_bodyweight = ?", isBodyweight)
	}

	// Get total count
	var total int64
	query.Count(&total)

	var exercises []models.Exercise
	if err := query.Offset(offset).Limit(params.Limit).Order("name ASC").Find(&exercises).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch exercises.")
		return
	}

	utils.PaginatedResponse(c, "Exercises retrieved successfully.", exercises, params.Page, params.Limit, int(total))
}

func GetExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, ok := utils.ParseUUID(c, params.ID, "exercise")
	if !ok {
		return
	}

	var exercise models.Exercise
	if err := database.DB.Where("id = ?", exerciseID).
		Preload("MuscleGroups.MuscleGroup").
		First(&exercise).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	utils.SuccessResponse(c, "Exercise retrieved successfully.", exercise)
}

func GetExerciseBySlug(c *gin.Context) {
	slug := c.Param("slug")
	// Validate slug format (alphanumeric with underscores/hyphens)
	if slug == "" {
		utils.BadRequestResponse(c, "Slug is required.", nil)
		return
	}

	var exercise models.Exercise
	if err := database.DB.Where("slug = ?", slug).
		Preload("MuscleGroups.MuscleGroup").
		Preload("Equipment.Equipment").
		First(&exercise).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	utils.SuccessResponse(c, "Exercise retrieved successfully.", exercise)
}

func UpdateExercise(c *gin.Context) {
	exerciseID, ok := utils.ParseUUIDParam(c, "id", "exercise")
	if !ok {
		return
	}

	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var exercise models.Exercise
	if err := database.DB.Where("id = ?", exerciseID).First(&exercise).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	// Update slug if name has changed
	if exercise.Name != req.Name {
		exercise.Slug = generateSlug(req.Name)
	}

	exercise.Name = req.Name
	exercise.Description = req.Description
	exercise.IsBodyweight = req.IsBodyweight
	exercise.Instructions = req.Instructions
	exercise.VideoURL = req.VideoURL

	if err := database.DB.Save(&exercise).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			utils.ConflictResponse(c, "Exercise with this name already exists.")
			return
		}
		utils.InternalServerErrorResponse(c, "Failed to update exercise.")
		return
	}

	// Load the exercise with muscle groups for response
	database.DB.Where("id = ?", exercise.ID).
		Preload("MuscleGroups.MuscleGroup").
		First(&exercise)

	utils.SuccessResponse(c, "Exercise updated successfully.", exercise)
}

func DeleteExercise(c *gin.Context) {
	exerciseID, ok := utils.ParseUUIDParam(c, "id", "exercise")
	if !ok {
		return
	}

	result := database.DB.Where("id = ?", exerciseID).Delete(&models.Exercise{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete exercise.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	utils.DeletedResponse(c, "Exercise deleted successfully.")
}
