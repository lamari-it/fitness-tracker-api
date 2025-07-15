package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateExerciseRequest struct {
	Name         string                    `json:"name" binding:"required"`
	Description  string                    `json:"description"`
	Equipment    string                    `json:"equipment"`
	IsBodyweight bool                      `json:"is_bodyweight"`
	Instructions string                    `json:"instructions"`
	VideoURL     string                    `json:"video_url"`
	MuscleGroups []MuscleGroupAssignment   `json:"muscle_groups,omitempty"`
}

type MuscleGroupAssignment struct {
	MuscleGroupID uuid.UUID `json:"muscle_group_id" binding:"required"`
	Primary       bool      `json:"primary"`
	Intensity     string    `json:"intensity"`
}

func CreateExercise(c *gin.Context) {
	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
		Name:         req.Name,
		Description:  req.Description,
		IsBodyweight: req.IsBodyweight,
		Instructions: req.Instructions,
		VideoURL:     req.VideoURL,
	}

	if err := tx.Create(&exercise).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Exercise with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only one muscle group can be set as primary"})
			return
		}

		for _, mgAssign := range req.MuscleGroups {
			// Verify muscle group exists
			var muscleGroup models.MuscleGroup
			if err := tx.Where("id = ?", mgAssign.MuscleGroupID).First(&muscleGroup).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid muscle group ID: " + mgAssign.MuscleGroupID.String()})
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid muscle group assignment"})
				return
			}

			if err := tx.Create(&assignment).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign muscle groups"})
				return
			}
		}
	}

	tx.Commit()

	// Load the exercise with muscle groups for response
	database.DB.Where("id = ?", exercise.ID).
		Preload("MuscleGroups.MuscleGroup").
		First(&exercise)

	c.JSON(http.StatusCreated, exercise)
}

func GetExercises(c *gin.Context) {
	search := c.Query("search")
	muscleGroupID := c.Query("muscle_group_id")
	equipment := c.Query("equipment")
	bodyweight := c.Query("bodyweight")
	primaryOnly := c.Query("primary_only")

	query := database.DB.Model(&models.Exercise{}).
		Preload("MuscleGroups.MuscleGroup")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if muscleGroupID != "" {
		// Filter by muscle group through the many-to-many relationship
		if primaryOnly == "true" {
			query = query.Joins("JOIN exercise_muscle_groups emg ON exercises.id = emg.exercise_id").
				Where("emg.muscle_group_id = ? AND emg.primary = true", muscleGroupID)
		} else {
			query = query.Joins("JOIN exercise_muscle_groups emg ON exercises.id = emg.exercise_id").
				Where("emg.muscle_group_id = ?", muscleGroupID)
		}
	}

	if equipment != "" {
		query = query.Where("equipment = ?", equipment)
	}

	if bodyweight != "" {
		isBodyweight := bodyweight == "true"
		query = query.Where("is_bodyweight = ?", isBodyweight)
	}

	var exercises []models.Exercise
	if err := query.Order("name ASC").Find(&exercises).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercises": exercises})
}

func GetExercise(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var exercise models.Exercise
	if err := database.DB.Where("id = ?", exerciseID).
		Preload("MuscleGroups.MuscleGroup").
		First(&exercise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

func UpdateExercise(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var req CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var exercise models.Exercise
	if err := database.DB.Where("id = ?", exerciseID).First(&exercise).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	exercise.Name = req.Name
	exercise.Description = req.Description
	exercise.IsBodyweight = req.IsBodyweight
	exercise.Instructions = req.Instructions
	exercise.VideoURL = req.VideoURL

	if err := database.DB.Save(&exercise).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{"error": "Exercise with this name already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update exercise"})
		return
	}

	// Load the exercise with muscle groups for response
	database.DB.Where("id = ?", exercise.ID).
		Preload("MuscleGroups.MuscleGroup").
		First(&exercise)

	c.JSON(http.StatusOK, exercise)
}

func DeleteExercise(c *gin.Context) {
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	result := database.DB.Where("id = ?", exerciseID).Delete(&models.Exercise{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete exercise"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercise deleted successfully"})
}