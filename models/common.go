package models

// WeightInput represents weight input from API requests
// Used for any weight field across the entire API
type WeightInput struct {
	WeightValue *float64 `json:"weight_value,omitempty"`
	WeightUnit  *string  `json:"weight_unit,omitempty" binding:"omitempty,oneof=kg lb"`
}

// WeightOutput represents weight output in API responses
// Always converted to user's preferred unit
type WeightOutput struct {
	WeightValue *float64 `json:"weight_value,omitempty"`
	WeightUnit  *string  `json:"weight_unit,omitempty"`
}
