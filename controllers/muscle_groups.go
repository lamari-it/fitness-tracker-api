package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateMuscleGroupRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Category    string `json:"category" binding:"omitempty,max=50"`
}

type UpdateMuscleGroupRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=100"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Category    string `json:"category" binding:"omitempty,max=50"`
}

type AssignMuscleGroupRequest struct {
	MuscleGroupID uuid.UUID `json:"muscle_group_id" binding:"required,uuid"`
	Primary       bool      `json:"primary"`
	Intensity     string    `json:"intensity" binding:"omitempty,oneof=low moderate high"`
}

func CreateMuscleGroup(c *gin.Context) {
	var req CreateMuscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	muscleGroup := models.MuscleGroup{
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
	}

	if err := muscleGroup.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed.", err.Error())
		return
	}

	if err := database.DB.Create(&muscleGroup).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create muscle group.")
		return
	}

	utils.CreatedResponse(c, "Muscle group created successfully.", muscleGroup.ToResponse())
}

func GetMuscleGroups(c *gin.Context) {
	var queryParams MuscleGroupQuery
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set default pagination values
	SetDefaultPagination(&queryParams.PaginationQuery)

	offset := (queryParams.Page - 1) * queryParams.Limit

	query := database.DB.Model(&models.MuscleGroup{})

	if queryParams.Search != "" {
		query = query.Where("name ILIKE ?", "%"+queryParams.Search+"%")
	}

	if queryParams.Category != "" {
		query = query.Where("category = ?", queryParams.Category)
	}

	var muscleGroups []models.MuscleGroup
	var total int64

	query.Count(&total)

	if err := query.Offset(offset).Limit(queryParams.Limit).Order("name ASC").Find(&muscleGroups).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch muscle groups.")
		return
	}

	var responses []models.MuscleGroupResponse
	for _, mg := range muscleGroups {
		responses = append(responses, mg.ToResponse())
	}

	utils.PaginatedResponse(c, "Muscle groups retrieved successfully.", responses, queryParams.Page, queryParams.Limit, int(total))
}

func GetMuscleGroup(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	muscleGroupID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid muscle group ID.", nil)
		return
	}

	var muscleGroup models.MuscleGroup
	if err := database.DB.Where("id = ?", muscleGroupID).
		Preload("ExerciseLinks.Exercise").
		First(&muscleGroup).Error; err != nil {
		utils.NotFoundResponse(c, "Muscle group not found.")
		return
	}

	response := models.MuscleGroupWithExercises{
		MuscleGroupResponse: muscleGroup.ToResponse(),
		ExerciseCount:       len(muscleGroup.ExerciseLinks),
	}

	for _, link := range muscleGroup.ExerciseLinks {
		exerciseResponse := models.ExerciseMuscleGroupResponse{
			ID:            link.ID,
			ExerciseID:    link.ExerciseID,
			MuscleGroupID: link.MuscleGroupID,
			Primary:       link.Primary,
			Intensity:     link.Intensity,
			MuscleGroup:   muscleGroup.ToResponse(),
		}
		response.Exercises = append(response.Exercises, exerciseResponse)
	}

	utils.SuccessResponse(c, "Muscle group retrieved successfully.", response)
}

func UpdateMuscleGroup(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	muscleGroupID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid muscle group ID.", nil)
		return
	}

	var req UpdateMuscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var muscleGroup models.MuscleGroup
	if err := database.DB.Where("id = ?", muscleGroupID).First(&muscleGroup).Error; err != nil {
		utils.NotFoundResponse(c, "Muscle group not found.")
		return
	}

	if req.Name != "" {
		muscleGroup.Name = req.Name
	}
	if req.Description != "" {
		muscleGroup.Description = req.Description
	}
	if req.Category != "" {
		muscleGroup.Category = req.Category
	}

	if err := muscleGroup.Validate(); err != nil {
		utils.BadRequestResponse(c, "Validation failed.", err.Error())
		return
	}

	if err := database.DB.Save(&muscleGroup).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update muscle group.")
		return
	}

	utils.SuccessResponse(c, "Muscle group updated successfully.", muscleGroup.ToResponse())
}

func DeleteMuscleGroup(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	muscleGroupID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid muscle group ID.", nil)
		return
	}

	// Check if muscle group is being used by any exercises
	var count int64
	database.DB.Model(&models.ExerciseMuscleGroup{}).Where("muscle_group_id = ?", muscleGroupID).Count(&count)
	if count > 0 {
		utils.ConflictResponse(c, "Cannot delete muscle group that is assigned to exercises.")
		return
	}

	result := database.DB.Where("id = ?", muscleGroupID).Delete(&models.MuscleGroup{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete muscle group.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Muscle group not found.")
		return
	}

	utils.DeletedResponse(c, "Muscle group deleted successfully.")
}

func AssignMuscleGroupToExercise(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise ID.", nil)
		return
	}

	var req AssignMuscleGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Check if exercise exists
	var exercise models.Exercise
	if err := database.DB.Where("id = ?", exerciseID).First(&exercise).Error; err != nil {
		utils.NotFoundResponse(c, "Exercise not found.")
		return
	}

	// Check if muscle group exists
	var muscleGroup models.MuscleGroup
	if err := database.DB.Where("id = ?", req.MuscleGroupID).First(&muscleGroup).Error; err != nil {
		utils.NotFoundResponse(c, "Muscle group not found.")
		return
	}

	// Check if assignment already exists
	var existing models.ExerciseMuscleGroup
	if err := database.DB.Where("exercise_id = ? AND muscle_group_id = ?", exerciseID, req.MuscleGroupID).First(&existing).Error; err == nil {
		utils.ConflictResponse(c, "Muscle group already assigned to this exercise.")
		return
	}

	// If this is being set as primary, unset any existing primary muscle groups
	if req.Primary {
		database.DB.Model(&models.ExerciseMuscleGroup{}).
			Where("exercise_id = ? AND primary = true", exerciseID).
			Update("primary", false)
	}

	assignment := models.ExerciseMuscleGroup{
		ExerciseID:    exerciseID,
		MuscleGroupID: req.MuscleGroupID,
		Primary:       req.Primary,
		Intensity:     req.Intensity,
	}

	if assignment.Intensity == "" {
		assignment.Intensity = "moderate"
	}

	if err := assignment.Validate(); err != nil {
		utils.BadRequestResponse(c, "Invalid assignment data.", err.Error())
		return
	}

	if err := database.DB.Create(&assignment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to assign muscle group to exercise.")
		return
	}

	// Load the muscle group for response
	database.DB.Where("id = ?", req.MuscleGroupID).First(&assignment.MuscleGroup)

	utils.CreatedResponse(c, "Muscle group assigned to exercise successfully.", assignment.ToResponse())
}

func RemoveMuscleGroupFromExercise(c *gin.Context) {
	var params struct {
		ID            string `uri:"id" binding:"required,uuid"`
		MuscleGroupID string `uri:"muscle_group_id" binding:"required,uuid"`
	}
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise ID.", nil)
		return
	}

	muscleGroupID, err := uuid.Parse(params.MuscleGroupID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid muscle group ID.", nil)
		return
	}

	result := database.DB.Where("exercise_id = ? AND muscle_group_id = ?", exerciseID, muscleGroupID).Delete(&models.ExerciseMuscleGroup{})
	if result.Error != nil {
		utils.InternalServerErrorResponse(c, "Failed to remove muscle group from exercise.")
		return
	}

	if result.RowsAffected == 0 {
		utils.NotFoundResponse(c, "Muscle group assignment not found.")
		return
	}

	utils.DeletedResponse(c, "Muscle group removed from exercise successfully.")
}

func GetExerciseMuscleGroups(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	exerciseID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid exercise ID.", nil)
		return
	}

	var assignments []models.ExerciseMuscleGroup
	if err := database.DB.Where("exercise_id = ?", exerciseID).
		Preload("MuscleGroup").
		Find(&assignments).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch exercise muscle groups.")
		return
	}

	var responses []models.ExerciseMuscleGroupResponse
	for _, assignment := range assignments {
		responses = append(responses, assignment.ToResponse())
	}

	utils.SuccessResponse(c, "Exercise muscle groups retrieved successfully.", responses)
}
