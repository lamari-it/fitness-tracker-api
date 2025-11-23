package utils

import "math"

// Weight conversion constants
// 1 pound = 0.45359237 kilograms (exact definition)
const (
	LbsToKgFactor = 0.45359237
	KgToLbsFactor = 2.20462262185 // 1 / 0.45359237
)

// ValidWeightUnits defines acceptable weight unit values
var ValidWeightUnits = map[string]bool{
	"kg":  true,
	"lb":  true,
	"lbs": true, // Alias for lb
}

// LbsToKg converts pounds to kilograms
func LbsToKg(lbs float64) float64 {
	return roundToDecimal(lbs*LbsToKgFactor, 2)
}

// KgToLbs converts kilograms to pounds
func KgToLbs(kg float64) float64 {
	return roundToDecimal(kg*KgToLbsFactor, 2)
}

// ConvertToKg converts a weight value from the given unit to kilograms
// Supported units: "kg", "lb", "lbs"
// If unit is empty or unrecognized, assumes kg
func ConvertToKg(weight float64, unit string) float64 {
	switch unit {
	case "lb", "lbs":
		return LbsToKg(weight)
	case "kg", "":
		return roundToDecimal(weight, 2)
	default:
		// Unknown unit, assume kg
		return roundToDecimal(weight, 2)
	}
}

// ConvertFromKg converts a weight value from kilograms to the target unit
// Supported units: "kg", "lb", "lbs"
// If unit is empty or unrecognized, returns kg
func ConvertFromKg(kg float64, unit string) float64 {
	switch unit {
	case "lb", "lbs":
		return KgToLbs(kg)
	case "kg", "":
		return roundToDecimal(kg, 2)
	default:
		// Unknown unit, return kg
		return roundToDecimal(kg, 2)
	}
}

// NormalizeWeightUnit normalizes weight unit strings to a standard format
// Returns "kg" or "lb"
func NormalizeWeightUnit(unit string) string {
	switch unit {
	case "lb", "lbs":
		return "lb"
	case "kg", "":
		return "kg"
	default:
		return "kg"
	}
}

// GetUserPreferredWeightUnit returns the user's preferred weight unit
// Now simply returns the value since it's already stored as kg/lb
// Returns "kg" as default if empty or invalid
func GetUserPreferredWeightUnit(preferredWeightUnit string) string {
	switch preferredWeightUnit {
	case "kg", "lb":
		return preferredWeightUnit
	default:
		return "kg"
	}
}

// roundToDecimal rounds a float to the specified number of decimal places
func roundToDecimal(value float64, decimals int) float64 {
	shift := math.Pow(10, float64(decimals))
	return math.Round(value*shift) / shift
}
