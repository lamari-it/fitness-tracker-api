package controllers

import (
	"fit-flow-api/database"
	"fit-flow-api/models"
	"fit-flow-api/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateWorkoutSession starts a new workout session
// If workout_id is provided and the workout has prescriptions, automatically creates
// session_blocks, session_exercises, and session_sets from the prescriptions
func CreateWorkoutSession(c *gin.Context) {
	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var req models.CreateWorkoutSessionRequest
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

	// Create the session
	session := models.WorkoutSession{
		UserID:      targetUserID,
		CreatedByID: &authUserID,
		WorkoutID:   req.WorkoutID,
		StartedAt:   startedAt,
		Notes:       "",
		Completed:   false,
	}

	if req.Notes != nil {
		session.Notes = *req.Notes
	}

	// Start transaction for creating session and related entities
	tx := database.DB.Begin()

	if err := tx.Create(&session).Error; err != nil {
		tx.Rollback()
		utils.InternalServerErrorResponse(c, "Failed to create workout session")
		return
	}

	// If workout_id provided, auto-create blocks/exercises/sets from prescriptions
	if req.WorkoutID != nil {
		if err := autoCreateSessionStructure(tx, session.ID, *req.WorkoutID); err != nil {
			tx.Rollback()
			utils.InternalServerErrorResponse(c, "Failed to create session structure from prescriptions")
			return
		}
	}

	tx.Commit()

	// Reload with full relationships
	database.DB.
		Preload("User").
		Preload("CreatedBy").
		Preload("Workout").
		Preload("SessionBlocks.SessionExercises.Exercise").
		Preload("SessionBlocks.SessionExercises.SessionSets.RPEValue").
		First(&session, "id = ?", session.ID)

	utils.CreatedResponse(c, "Workout session created successfully", models.BuildSessionResponse(session, "kg"))
}

// autoCreateSessionStructure creates session blocks, exercises, and sets from workout prescriptions
func autoCreateSessionStructure(tx *gorm.DB, sessionID uuid.UUID, workoutID uuid.UUID) error {
	// Get all prescriptions for this workout, ordered by group_order and exercise_order
	var prescriptions []models.WorkoutPrescription
	if err := tx.
		Where("workout_id = ?", workoutID).
		Order("group_order ASC, exercise_order ASC").
		Find(&prescriptions).Error; err != nil {
		return err
	}

	if len(prescriptions) == 0 {
		return nil // No prescriptions, just create empty session
	}

	// Group prescriptions by GroupID
	groupMap := make(map[uuid.UUID][]models.WorkoutPrescription)
	groupOrder := make([]uuid.UUID, 0)

	for _, p := range prescriptions {
		if _, exists := groupMap[p.GroupID]; !exists {
			groupOrder = append(groupOrder, p.GroupID)
		}
		groupMap[p.GroupID] = append(groupMap[p.GroupID], p)
	}

	// Create session blocks for each prescription group
	for blockOrder, groupID := range groupOrder {
		groupPrescriptions := groupMap[groupID]
		firstPrescription := groupPrescriptions[0]

		// Create session block
		block := models.SessionBlock{
			SessionID:  sessionID,
			GroupID:    groupID,
			BlockOrder: blockOrder + 1,
			Skipped:    false,
		}

		if err := tx.Create(&block).Error; err != nil {
			return err
		}

		// Create session exercises for each prescription in the group
		for _, prescription := range groupPrescriptions {
			sessionExercise := models.SessionExercise{
				SessionBlockID: block.ID,
				PrescriptionID: &prescription.ID,
				ExerciseID:     prescription.ExerciseID,
				ExerciseOrder:  prescription.ExerciseOrder,
				Skipped:        false,
				Notes:          "",
			}

			if err := tx.Create(&sessionExercise).Error; err != nil {
				return err
			}

			// Create session sets based on prescription
			numSets := 1
			if prescription.Sets != nil && *prescription.Sets > 0 {
				numSets = *prescription.Sets
			}

			// For types with group_rounds, multiply by rounds
			if firstPrescription.GroupRounds != nil && *firstPrescription.GroupRounds > 1 {
				// For circuits/supersets, we create sets based on rounds
				// But for now, keep it simple - just use the sets value
			}

			for setNum := 1; setNum <= numSets; setNum++ {
				set := models.SessionSet{
					SessionExerciseID: sessionExercise.ID,
					SetNumber:         setNum,
					Completed:         false,
					ActualWeightKg:    prescription.TargetWeightKg, // Pre-fill with target
					RPEValueID:        prescription.RPEValueID,
					WasFailure:        false,
					Notes:             "",
				}

				// Pre-fill reps or duration based on prescription
				if prescription.Reps != nil {
					set.ActualReps = prescription.Reps
				}
				if prescription.HoldSeconds != nil {
					set.ActualDurationSeconds = prescription.HoldSeconds
				}

				if err := tx.Create(&set).Error; err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GetWorkoutSession retrieves a single workout session
func GetWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, ok := utils.ParseUUID(c, params.ID, "workout session")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	// Get weight unit preference from query param, default to kg
	weightUnit := c.DefaultQuery("unit", "kg")
	if weightUnit != "kg" && weightUnit != "lb" {
		weightUnit = "kg"
	}

	var session models.WorkoutSession
	if err := database.DB.
		Preload("User").
		Preload("CreatedBy").
		Preload("Workout").
		Preload("SessionBlocks", func(db *gorm.DB) *gorm.DB {
			return db.Order("block_order ASC")
		}).
		Preload("SessionBlocks.SessionExercises", func(db *gorm.DB) *gorm.DB {
			return db.Order("exercise_order ASC")
		}).
		Preload("SessionBlocks.SessionExercises.Exercise").
		Preload("SessionBlocks.SessionExercises.SessionSets", func(db *gorm.DB) *gorm.DB {
			return db.Order("set_number ASC")
		}).
		Preload("SessionBlocks.SessionExercises.SessionSets.RPEValue").
		First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization: user owns the session OR trainer created it
	if !isAuthorizedForSession(session, authUserID) {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	utils.SuccessResponse(c, "Workout session retrieved successfully", models.BuildSessionResponse(session, weightUnit))
}

// GetWorkoutSessions lists workout sessions for the authenticated user
func GetWorkoutSessions(c *gin.Context) {
	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	// Check if querying for a specific client (trainer view)
	clientIDStr := c.Query("client_id")
	var targetUserID uuid.UUID

	if clientIDStr != "" {
		clientID, ok := utils.ParseUUID(c, clientIDStr, "client")
		if !ok {
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

	// Get weight unit preference from query param, default to kg
	weightUnit := c.DefaultQuery("unit", "kg")
	if weightUnit != "kg" && weightUnit != "lb" {
		weightUnit = "kg"
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
		responses[i] = models.BuildSessionResponse(session, weightUnit)
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

	sessionID, ok := utils.ParseUUID(c, params.ID, "workout session")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to end this session")
		return
	}

	var req models.EndWorkoutSessionRequest
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
	session.Completed = true

	// Calculate duration
	duration := int(endedAt.Sub(session.StartedAt).Seconds())
	session.DurationSeconds = &duration

	if req.Notes != nil {
		session.Notes = *req.Notes
	}

	if req.PerceivedIntensity != nil {
		session.PerceivedIntensity = req.PerceivedIntensity
	}

	if err := database.DB.Save(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to end workout session")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("User").
		Preload("CreatedBy").
		Preload("Workout").
		First(&session, "id = ?", session.ID)

	utils.SuccessResponse(c, "Workout session ended successfully", models.BuildSessionResponse(session, "kg"))
}

// UpdateWorkoutSession updates session notes and perceived intensity
func UpdateWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, ok := utils.ParseUUID(c, params.ID, "workout session")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to update this session")
		return
	}

	var req models.UpdateWorkoutSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	if req.Notes != nil {
		session.Notes = *req.Notes
	}

	if req.PerceivedIntensity != nil {
		session.PerceivedIntensity = req.PerceivedIntensity
	}

	if err := database.DB.Save(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to update workout session")
		return
	}

	// Reload with relationships
	database.DB.
		Preload("User").
		Preload("CreatedBy").
		Preload("Workout").
		First(&session, "id = ?", session.ID)

	utils.SuccessResponse(c, "Workout session updated successfully", models.BuildSessionResponse(session, "kg"))
}

// DeleteWorkoutSession deletes a workout session and all related data
func DeleteWorkoutSession(c *gin.Context) {
	var params IDParam
	if err := c.ShouldBindUri(&params); err != nil {
		utils.HandleBindingError(c, err)
		return
	}

	sessionID, ok := utils.ParseUUID(c, params.ID, "workout session")
	if !ok {
		return
	}

	authUserID, ok := utils.GetAuthUserID(c)
	if !ok {
		return
	}

	var session models.WorkoutSession
	if err := database.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		utils.NotFoundResponse(c, "Workout session not found")
		return
	}

	// Authorization
	if !isAuthorizedForSession(session, authUserID) {
		utils.ForbiddenResponse(c, "Not authorized to delete this session")
		return
	}

	// Delete cascades through foreign keys in the database
	if err := database.DB.Delete(&session).Error; err != nil {
		utils.InternalServerErrorResponse(c, "Failed to delete workout session")
		return
	}

	utils.NoContentResponse(c)
}
