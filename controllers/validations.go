package controllers

// PaginationQuery represents common pagination parameters with validation
type PaginationQuery struct {
	Page  int `form:"page" validate:"min=1" binding:"omitempty,min=1"`
	Limit int `form:"limit" validate:"min=1,max=100" binding:"omitempty,min=1,max=100"`
}

// ExerciseQuery represents query parameters for exercise endpoints
type ExerciseQuery struct {
	PaginationQuery
	Search        string `form:"search" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	MuscleGroupID string `form:"muscle_group_id" validate:"omitempty,uuid" binding:"omitempty,uuid"`
	Equipment     string `form:"equipment" validate:"omitempty,max=50" binding:"omitempty,max=50"`
	Bodyweight    string `form:"bodyweight" validate:"omitempty,oneof=true false" binding:"omitempty,oneof=true false"`
	PrimaryOnly   string `form:"primary_only" validate:"omitempty,oneof=true false" binding:"omitempty,oneof=true false"`
	IsFavorited   string `form:"is_favorited" validate:"omitempty,oneof=true false" binding:"omitempty,oneof=true false"`
}

// WorkoutQuery represents query parameters for workout endpoints
// Supports multiple muscle_group_id and exercise_id values via repeated query params
type WorkoutQuery struct {
	PaginationQuery
	Search         string   `form:"search" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	MuscleGroupIDs []string `form:"muscle_group_id"`
	ExerciseIDs    []string `form:"exercise_id"`
	Mode           string   `form:"mode" validate:"omitempty,oneof=and or" binding:"omitempty,oneof=and or"`
	IsFavorited    string   `form:"is_favorited" validate:"omitempty,oneof=true false" binding:"omitempty,oneof=true false"`
}

// GetFilterMode returns the filter mode, defaulting to "or"
func (q *WorkoutQuery) GetFilterMode() string {
	if q.Mode == "" {
		return "or"
	}
	return q.Mode
}

// EquipmentQuery represents query parameters for equipment endpoints
type EquipmentQuery struct {
	PaginationQuery
	Search   string `form:"search" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	Category string `form:"category" validate:"omitempty,max=50" binding:"omitempty,max=50"`
}

// MuscleGroupQuery represents query parameters for muscle group endpoints
type MuscleGroupQuery struct {
	PaginationQuery
	Search   string `form:"search" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	Category string `form:"category" validate:"omitempty,max=50" binding:"omitempty,max=50"`
}

// FitnessGoalQuery represents query parameters for fitness goal endpoints
type FitnessGoalQuery struct {
	PaginationQuery
	Category string `form:"category" validate:"omitempty,max=50" binding:"omitempty,max=50"`
}

// TranslationQuery represents query parameters for translation endpoints
type TranslationQuery struct {
	PaginationQuery
	ResourceType string `form:"resource_type" validate:"omitempty,max=50" binding:"omitempty,max=50"`
	ResourceID   string `form:"resource_id" validate:"omitempty,uuid" binding:"omitempty,uuid"`
	Language     string `form:"language" validate:"omitempty,len=2" binding:"omitempty,len=2"`
}

// EnrollmentQuery represents query parameters for enrollment endpoints
type EnrollmentQuery struct {
	PaginationQuery
	Status string `form:"status" validate:"omitempty,oneof=active paused completed cancelled" binding:"omitempty,oneof=active paused completed cancelled"`
}

// UserEquipmentQuery represents query parameters for user equipment endpoints
type UserEquipmentQuery struct {
	PaginationQuery
	LocationType string `form:"location_type" validate:"omitempty,oneof=home gym" binding:"omitempty,oneof=home gym"`
}

// TrainerQuery represents query parameters for trainer endpoints
type TrainerQuery struct {
	PaginationQuery
	Search    string  `form:"search" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	Specialty string  `form:"specialty" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	Location  string  `form:"location" validate:"omitempty,max=100" binding:"omitempty,max=100"`
	MinRating float64 `form:"min_rating" validate:"omitempty,min=0,max=5" binding:"omitempty,min=0,max=5"`
	SortBy    string  `form:"sort_by" validate:"omitempty,oneof=rating rate recent" binding:"omitempty,oneof=rating rate recent"`
}

// IDParam represents common UUID path parameters
type IDParam struct {
	ID string `uri:"id" binding:"required,uuid"`
}

// SlugParam represents slug path parameters
type SlugParam struct {
	Slug string `uri:"slug" binding:"required,alphanum"`
}

// SetDefaultPagination sets default values for pagination if not provided
func SetDefaultPagination(query *PaginationQuery) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}
}

// GetOffset calculates the database offset for pagination.
func (p *PaginationQuery) GetOffset() int {
	return (p.Page - 1) * p.Limit
}
