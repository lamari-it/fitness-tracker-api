package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request DTOs
type CreateWorkoutSessionRequest struct {
	UserID    *uuid.UUID `json:"user_id"`    // Optional: if not provided, uses authenticated user
	WorkoutID *uuid.UUID `json:"workout_id"` // Optional: for free-form workouts
	StartedAt *time.Time `json:"started_at"` // Optional: defaults to now
	Notes     string     `json:"notes"`
}

type EndWorkoutSessionRequest struct {
	EndedAt *time.Time `json:"ended_at"` // Optional: defaults to now
	Notes   string     `json:"notes"`
}

type UpdateWorkoutSessionRequest struct {
	Notes string `json:"notes"`
}

// CreateWorkoutSession starts a new workout session
func CreateWorkoutSession(c *gin.Context) {
	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var req CreateWorkoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Determine target user (who the session is for)
	targetUserID := authUserID
	if req.UserID != nil && *req.UserID != authUserID {
		// Creating session for someone else - verify trainer-client relationship
		var link models.TrainerClientLink
		if err := database.DB.Where(
			"trainer_id = ? AND client_id = ? AND status = ?",
			authUserID, *req.UserID, "active",
		).First(&link).Error; err != nil {
			utils.ForbiddenResponse(c, "Not authorized to create sessions for this user")
			return
		}
		targetUserID = *req.UserID
	}

	// Set start time
	startedAt := time.Now()
	if req.StartedAt != nil {
		startedAt = *req.StartedAt
	}

	session := models.WorkoutSession{
		UserID:      targetUserID,
		CreatedByID: &authUserID,
		WorkoutID:   req.WorkoutID,
		StartedAt:   startedAt,
		Notes:       req.Notes,
	}

	if err := database.DB.Create(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to create workout session")
		return
	}

	// Reload with relationships
	database.DB.Preload("User").Preload("CreatedBy").Preload("Workout").First(&session, "id = ?", session.ID)

	utils.CreatedResponse(c, "Workout session created successfully", session.ToResponse())
}

// GetWorkoutSession retrieves a single workout session
func GetWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var session models.WorkoutSession
	if err := database.DB.
		Preload("User").
		Preload("CreatedBy").
		Preload("Workout").
		Preload("ExerciseLogs.Exercise").
		Preload("ExerciseLogs.SetGroup").
		Preload("ExerciseLogs.SetLogs").
		First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization: user owns the session OR trainer created it
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID

	if !isOwner && !isCreator {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	utils.SuccessResponse(c, "Workout session retrieved successfully", session.ToResponse())
}

// GetWorkoutSessions lists workout sessions for the authenticated user
func GetWorkoutSessions(c *gin.Context) {
	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	// Check if querying for a specific client (trainer view)
	clientIDStr := c.Query("client_id")
	var targetUserID uuid.UUID

	if clientIDStr != "" {
		clientID, err := uuid.Parse(clientIDStr)
		if err != nil {
			utils.BadRequestResponse(c, "Invalid client_id format", nil)
			return
		}

		// Verify trainer-client relationship
		var link models.TrainerClientLink
		if err := database.DB.Where(
			"trainer_id = ? AND client_id = ? AND status = ?",
			authUserID, clientID, "active",
		).First(&link).Error; err != nil {
			utils.ForbiddenResponse(c, "Not authorized to view sessions for this user")
			return
		}

		targetUserID = clientID
	} else {
		targetUserID = authUserID
	}

	var sessions []models.WorkoutSession
	query := database.DB.
		Preload("User").
		Preload("CreatedBy").
		Preload("Workout").
		Where("user_id = ?", targetUserID)

	// If trainer is viewing client sessions, only show sessions they created
	if clientIDStr != "" {
		query = query.Where("created_by_id = ?", authUserID)
	}

	if err := query.Order("started_at DESC").Find(&sessions).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to retrieve workout sessions")
		return
	}

	responses := make([]models.WorkoutSessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = session.ToResponse()
	}

	utils.SuccessResponse(c, "Workout sessions retrieved successfully", responses)
}

// EndWorkoutSession marks a workout session as complete
func EndWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization: user owns the session OR trainer created it
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID

	if !isOwner && !isCreator {
		utils.ForbiddenResponse(c, "Not authorized to end this session")
		return
	}

	var req EndWorkoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	// Set end time
	endedAt := time.Now()
	if req.EndedAt != nil {
		endedAt = *req.EndedAt
	}

	session.EndedAt = &endedAt
	if req.Notes != "" {
		session.Notes = req.Notes
	}

	if err := database.DB.Save(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to end workout session")
		return
	}

	// Reload with relationships
	database.DB.Preload("User").Preload("CreatedBy").Preload("Workout").First(&session, "id = ?", session.ID)

	utils.SuccessResponse(c, "Workout session ended successfully", session.ToResponse())
}

// UpdateWorkoutSession updates session notes
func UpdateWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization: user owns the session OR trainer created it
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID

	if !isOwner && !isCreator {
		utils.ForbiddenResponse(c, "Not authorized to update this session")
		return
	}

	var req UpdateWorkoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	session.Notes = req.Notes

	if err := database.DB.Save(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update workout session")
		return
	}

	// Reload with relationships
	database.DB.Preload("User").Preload("CreatedBy").Preload("Workout").First(&session, "id = ?", session.ID)

	utils.SuccessResponse(c, "Workout session updated successfully", session.ToResponse())
}

// DeleteWorkoutSession deletes a workout session
func DeleteWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, err := uuid.Parse(params.ID)
	if err != nil {
		utils.BadRequestResponse(c, "Invalid UUID format", nil)
		return
	}

	authUserIDVal, exists := c.Get("user_id")
	if !exists {
		utils.UnauthorizedResponse(c, "User not authenticated")
		return
	}

	authUserID, ok := authUserIDVal.(uuid.UUID)
	if !ok {
		utils.InternalServerErrorResponse(c, "Invalid user ID type")
		return
	}

	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization: user owns the session OR trainer created it
	isOwner := session.UserID == authUserID
	isCreator := session.CreatedByID != nil && *session.CreatedByID == authUserID

	if !isOwner && !isCreator {
		utils.ForbiddenResponse(c, "Not authorized to delete this session")
		return
	}

	if err := database.DB.Delete(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete workout session")
		return
	}

	utils.NoContentResponse(c)
}
