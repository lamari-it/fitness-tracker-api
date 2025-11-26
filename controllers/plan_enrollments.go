package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type EnrollInPlanRequest struct {
	PlanID            uuid.UUID `json:"plan_id" binding:"required"`
	StartDate         string    `json:"start_date" binding:"required"`
	DaysPerWeek       int       `json:"days_per_week" binding:"required,min=1,max=7"`
	ScheduleMode      string    `json:"schedule_mode" binding:"omitempty,oneof=rolling calendar"` // rolling or calendar
	PreferredWeekdays []int32   `json:"preferred_weekdays" binding:"omitempty,dive,min=0,max=6"`  // Only for calendar mode
}

type UpdateEnrollmentRequest struct {
	DaysPerWeek       int     `json:"days_per_week" binding:"omitempty,min=1,max=7"`
	ScheduleMode      string  `json:"schedule_mode" binding:"omitempty,oneof=rolling calendar"`
	PreferredWeekdays []int32 `json:"preferred_weekdays" binding:"omitempty,dive,min=0,max=6"`
	Status            string  `json:"status" binding:"omitempty,oneof=active paused completed cancelled"`
}

func EnrollInPlan(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req EnrollInPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Parse start date
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		validationErrors := utils.ValidationErrors{
			"start_date": []string{"Invalid date format. Use YYYY-MM-DD."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Check if plan exists
	var plan models.WorkoutPlan
	if err := database.DB.Where("id = ?", req.PlanID).First(&plan).Error; err != nil {
		utils.NotFoundResponse(c, "Workout plan not found.")
		return
	}

	// Check if user is already enrolled in this plan
	var existingEnrollment models.PlanEnrollment
	if err := database.DB.Where("user_id = ? AND plan_id = ? AND status = ?", userID, req.PlanID, "active").First(&existingEnrollment).Error; err == nil {
		utils.BadRequestResponse(c, "You are already enrolled in this plan.", nil)
		return
	}

	// Set default schedule mode if not provided
	scheduleMode := req.ScheduleMode
	if scheduleMode == "" {
		scheduleMode = "rolling"
	}

	// Validate schedule mode
	if scheduleMode != "rolling" && scheduleMode != "calendar" {
		validationErrors := utils.ValidationErrors{
			"schedule_mode": []string{"Schedule mode must be 'rolling' or 'calendar'."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	// Validate preferred weekdays for calendar mode
	if scheduleMode == "calendar" && len(req.PreferredWeekdays) != req.DaysPerWeek {
		validationErrors := utils.ValidationErrors{
			"preferred_weekdays": []string{"Number of preferred weekdays must match days per week for calendar mode."},
		}
		utils.ValidationErrorResponse(c, validationErrors)
		return
	}

	enrollment := models.PlanEnrollment{
		PlanID:            req.PlanID,
		UserID:            userID,
		StartDate:         startDate,
		DaysPerWeek:       req.DaysPerWeek,
		ScheduleMode:      scheduleMode,
		PreferredWeekdays: pq.Int32Array(req.PreferredWeekdays),
		Status:            "active",
		CurrentIndex:      0,
	}

	if err := database.DB.Create(&enrollment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to enroll in workout plan.")
		return
	}

	utils.CreatedResponse(c, "Successfully enrolled in workout plan.", enrollment)
}

func GetUserEnrollments(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	var enrollments []models.PlanEnrollment
	var total int64

	query := database.DB.Model(&models.PlanEnrollment{}).Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to count enrollments.")
		return
	}

	if err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&enrollments).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to fetch enrollments.")
		return
	}

	utils.PaginatedResponse(c, "Enrollments fetched successfully.", enrollments, page, limit, int(total))
}

func GetEnrollment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	enrollmentID, ok := utils.ParseUUID(c, c.Param("id"), "enrollment")
	if !ok {
		return
	}

	var enrollment models.PlanEnrollment
	if err := database.DB.Where("id = ? AND user_id = ?", enrollmentID, userID).First(&enrollment).Error; err != nil {
		utils.NotFoundResponse(c, "Enrollment not found.")
		return
	}

	utils.SuccessResponse(c, "Enrollment fetched successfully.", enrollment)
}

func UpdateEnrollment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	enrollmentID, ok := utils.ParseUUID(c, c.Param("id"), "enrollment")
	if !ok {
		return
	}

	var req UpdateEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	var enrollment models.PlanEnrollment
	if err := database.DB.Where("id = ? AND user_id = ?", enrollmentID, userID).First(&enrollment).Error; err != nil {
		utils.NotFoundResponse(c, "Enrollment not found.")
		return
	}

	updates := map[string]interface{}{}

	if req.DaysPerWeek > 0 && req.DaysPerWeek <= 7 {
		updates["days_per_week"] = req.DaysPerWeek
	}

	if req.ScheduleMode != "" {
		if req.ScheduleMode != "rolling" && req.ScheduleMode != "calendar" {
			validationErrors := utils.ValidationErrors{
				"schedule_mode": []string{"Schedule mode must be 'rolling' or 'calendar'."},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		updates["schedule_mode"] = req.ScheduleMode

		if req.ScheduleMode == "calendar" && len(req.PreferredWeekdays) > 0 {
			updates["preferred_weekdays"] = pq.Int32Array(req.PreferredWeekdays)
		}
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "paused" && req.Status != "completed" {
			validationErrors := utils.ValidationErrors{
				"status": []string{"Status must be 'active', 'paused', or 'completed'."},
			}
			utils.ValidationErrorResponse(c, validationErrors)
			return
		}
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&enrollment).Updates(updates).Error; err != nil {
			utils.InternalServerErrorResponse(c, "Failed to update enrollment.")
			return
		}
	}

	database.DB.Where("id = ?", enrollmentID).First(&enrollment)

	utils.SuccessResponse(c, "Enrollment updated successfully.", enrollment)
}

func CancelEnrollment(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	enrollmentID, ok := utils.ParseUUID(c, c.Param("id"), "enrollment")
	if !ok {
		return
	}

	var enrollment models.PlanEnrollment
	if err := database.DB.Where("id = ? AND user_id = ?", enrollmentID, userID).First(&enrollment).Error; err != nil {
		utils.NotFoundResponse(c, "Enrollment not found.")
		return
	}

	if enrollment.Status != "active" {
		utils.BadRequestResponse(c, "Only active enrollments can be cancelled.", nil)
		return
	}

	enrollment.Status = "cancelled"
	if err := database.DB.Save(&enrollment).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to cancel enrollment.")
		return
	}

	utils.SuccessResponse(c, "Enrollment cancelled successfully.", enrollment)
}
