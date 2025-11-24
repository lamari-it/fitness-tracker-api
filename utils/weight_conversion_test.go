package utils

import (
	"math"
	"testing"
)

func TestLbsToKg(t *testing.T) {
	tests := []struct {
		name     string
		lbs      float64
		expected float64
	}{
		{"Zero", 0, 0},
		{"One pound", 1, 0.45},
		{"100 pounds", 100, 45.36},
		{"135 pounds (common plate weight)", 135, 61.23},
		{"225 pounds (2 plates)", 225, 102.06},
		{"Small value", 0.5, 0.23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LbsToKg(tt.lbs)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("LbsToKg(%v) = %v, want %v", tt.lbs, result, tt.expected)
			}
		})
	}
}

func TestKgToLbs(t *testing.T) {
	tests := []struct {
		name     string
		kg       float64
		expected float64
	}{
		{"Zero", 0, 0},
		{"One kilogram", 1, 2.2},
		{"50 kg", 50, 110.23},
		{"100 kg", 100, 220.46},
		{"Small value", 0.5, 1.1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KgToLbs(tt.kg)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("KgToLbs(%v) = %v, want %v", tt.kg, result, tt.expected)
			}
		})
	}
}

func TestConvertToKg(t *testing.T) {
	tests := []struct {
		name     string
		weight   float64
		unit     string
		expected float64
	}{
		{"Kg to kg", 100, "kg", 100},
		{"Lb to kg", 100, "lb", 45.36},
		{"Lbs to kg", 100, "lbs", 45.36},
		{"Empty unit defaults to kg", 100, "", 100},
		{"Unknown unit defaults to kg", 100, "oz", 100},
		{"Zero weight", 0, "lb", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToKg(tt.weight, tt.unit)
			if math.Abs(result-tt.expected) > 0.01 {
				t.Errorf("ConvertToKg(%v, %v) = %v, want %v", tt.weight, tt.unit, result, tt.expected)
			}
		})
	}
}

func TestConvertFromKg(t *testing.T) {
	tests := []struct {
		name     string
		kg       float64
		unit     string
		expected float64
	}{
		{"Kg to kg", 100, "kg", 100},
		{"Kg to lb", 45.36, "lb", 100.0},
		{"Kg to lbs", 45.36, "lbs", 100.0},
		{"Empty unit defaults to kg", 100, "", 100},
		{"Unknown unit defaults to kg", 100, "oz", 100},
		{"Zero weight", 0, "lb", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertFromKg(tt.kg, tt.unit)
			if math.Abs(result-tt.expected) > 0.1 {
				t.Errorf("ConvertFromKg(%v, %v) = %v, want %v", tt.kg, tt.unit, result, tt.expected)
			}
		})
	}
}

func TestNormalizeWeightUnit(t *testing.T) {
	tests := []struct {
		name     string
		unit     string
		expected string
	}{
		{"kg stays kg", "kg", "kg"},
		{"lb stays lb", "lb", "lb"},
		{"lbs becomes lb", "lbs", "lb"},
		{"Empty becomes kg", "", "kg"},
		{"Unknown becomes kg", "oz", "kg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NormalizeWeightUnit(tt.unit)
			if result != tt.expected {
				t.Errorf("NormalizeWeightUnit(%v) = %v, want %v", tt.unit, result, tt.expected)
			}
		})
	}
}

func TestGetUserPreferredWeightUnit(t *testing.T) {
	tests := []struct {
		name         string
		weightUnit   string
		expectedUnit string
	}{
		{"kg returns kg", "kg", "kg"},
		{"lb returns lb", "lb", "lb"},
		{"Empty returns kg", "", "kg"},
		{"Unknown returns kg", "unknown", "kg"},
		{"Invalid returns kg", "lbs", "kg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserPreferredWeightUnit(tt.weightUnit)
			if result != tt.expectedUnit {
				t.Errorf("GetUserPreferredWeightUnit(%v) = %v, want %v", tt.weightUnit, result, tt.expectedUnit)
			}
		})
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test that converting kg -> lb -> kg gives approximately the same value
	tests := []float64{50, 100, 61.23, 102.06}

	for _, kg := range tests {
		lbs := KgToLbs(kg)
		backToKg := LbsToKg(lbs)

		if math.Abs(backToKg-kg) > 0.02 {
			t.Errorf("Round trip failed for %v kg: got %v (via %v lbs)", kg, backToKg, lbs)
		}
	}
}
