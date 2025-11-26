package controllers

import (
	"lamari-fit-api/database"
	"lamari-fit-api/models"
	"lamari-fit-api/utils"
	"math"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateWeightLog creates a new weight log entry
func CreateWeightLog(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req models.CreateWeightLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	weightLog := models.WeightLog{
		UserID:   userID,
		WeightKg: req.WeightKg,
		Notes:    req.Notes,
	}

	if err := database.DB.Create(&weightLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create weight log")
		return
	}

	// Also update current weight in fitness profile if it exists
	database.DB.Model(&models.UserFitnessProfile{}).
		Where("user_id = ?", userID).
		Update("current_weight_kg", req.WeightKg)

	utils.CreatedResponse(c, "Weight logged successfully", weightLog.ToResponse())
}

// GetWeightLogs retrieves weight history with pagination and date filtering
func GetWeightLogs(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if limit > 100 {
		limit = 100
	}
	offset := (page - 1) * limit

	query := database.DB.Where("user_id = ?", userID)

	// Apply date filters
	if startDate != "" {
		if parsed, err := time.Parse("2006-01-02", startDate); err == nil {
			query = query.Where("created_at >= ?", parsed)
		}
	}
	if endDate != "" {
		if parsed, err := time.Parse("2006-01-02", endDate); err == nil {
			// Add one day to include the entire end date
			query = query.Where("created_at < ?", parsed.AddDate(0, 0, 1))
		}
	}

	// Get total count
	var total int64
	query.Model(&models.WeightLog{}).Count(&total)

	// Get paginated results
	var logs []models.WeightLog
	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve weight logs")
		return
	}

	// Convert to response format
	responses := make([]models.WeightLogResponse, len(logs))
	for i, log := range logs {
		responses[i] = log.ToResponse()
	}

	utils.PaginatedResponse(c, "Weight logs retrieved successfully", responses, page, limit, int(total))
}

// GetWeightLog retrieves a specific weight log entry
func GetWeightLog(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	logID, ok := utils.ParseUUID(c, c.Param("id"), "weight log")
	if !ok {
		return
	}

	var weightLog models.WeightLog
	if err := database.DB.Where("id = ? AND user_id = ?", logID, userID).First(&weightLog).Error; err != nil {
		utils.NotFoundResponse(c, "Weight log not found")
		return
	}

	utils.SuccessResponse(c, "Weight log retrieved successfully", weightLog.ToResponse())
}

// UpdateWeightLog updates a weight log entry (only within 24 hours of creation)
func UpdateWeightLog(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	logID, ok := utils.ParseUUID(c, c.Param("id"), "weight log")
	if !ok {
		return
	}

	var weightLog models.WeightLog
	if err := database.DB.Where("id = ? AND user_id = ?", logID, userID).First(&weightLog).Error; err != nil {
		utils.NotFoundResponse(c, "Weight log not found")
		return
	}

	// Check if within 24 hours of creation
	if time.Since(weightLog.CreatedAt) > 24*time.Hour {
		utils.BadRequestResponse(c, "Weight logs can only be updated within 24 hours of creation", nil)
		return
	}

	var req models.UpdateWeightLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Update fields if provided
	if req.WeightKg > 0 {
		weightLog.WeightKg = req.WeightKg
	}
	if req.Notes != "" {
		weightLog.Notes = req.Notes
	}

	if err := database.DB.Save(&weightLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update weight log")
		return
	}

	utils.SuccessResponse(c, "Weight log updated successfully", weightLog.ToResponse())
}

// DeleteWeightLog deletes a weight log entry (soft delete)
func DeleteWeightLog(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	logID, ok := utils.ParseUUID(c, c.Param("id"), "weight log")
	if !ok {
		return
	}

	var weightLog models.WeightLog
	if err := database.DB.Where("id = ? AND user_id = ?", logID, userID).First(&weightLog).Error; err != nil {
		utils.NotFoundResponse(c, "Weight log not found")
		return
	}

	if err := database.DB.Delete(&weightLog).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete weight log")
		return
	}

	utils.NoContentResponse(c)
}

// GetWeightStats retrieves weight statistics for a given period
func GetWeightStats(c *gin.Context) {
	userID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	// Parse days parameter (default 30 days)
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	if days > 365 {
		days = 365
	}

	startDate := time.Now().AddDate(0, 0, -days)

	var logs []models.WeightLog
	if err := database.DB.Where("user_id = ? AND created_at >= ?", userID, startDate).
		Order("created_at ASC").Find(&logs).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve weight stats")
		return
	}

	if len(logs) == 0 {
		utils.NotFoundResponse(c, "No weight logs found for the specified period")
		return
	}

	// Calculate statistics
	var minWeight, maxWeight, sumWeight float64
	minWeight = logs[0].WeightKg
	maxWeight = logs[0].WeightKg

	for _, log := range logs {
		if log.WeightKg < minWeight {
			minWeight = log.WeightKg
		}
		if log.WeightKg > maxWeight {
			maxWeight = log.WeightKg
		}
		sumWeight += log.WeightKg
	}

	avgWeight := sumWeight / float64(len(logs))
	startWeight := logs[0].WeightKg
	latestWeight := logs[len(logs)-1].WeightKg
	weightChange := latestWeight - startWeight

	stats := models.WeightStatsResponse{
		TotalEntries: len(logs),
		LatestWeight: math.Round(latestWeight*100) / 100,
		MinWeight:    math.Round(minWeight*100) / 100,
		MaxWeight:    math.Round(maxWeight*100) / 100,
		AvgWeight:    math.Round(avgWeight*100) / 100,
		StartWeight:  math.Round(startWeight*100) / 100,
		WeightChange: math.Round(weightChange*100) / 100,
		PeriodDays:   days,
		StartDate:    logs[0].CreatedAt,
		EndDate:      logs[len(logs)-1].CreatedAt,
	}

	utils.SuccessResponse(c, "Weight statistics retrieved successfully", stats)
}
