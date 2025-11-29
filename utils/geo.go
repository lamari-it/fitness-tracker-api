package utils

import (
	"fmt"
	"math"
)

// EarthRadiusKm is the Earth's radius in kilometers
const EarthRadiusKm = 6371.0

// DefaultRadiusKm is the default search radius when not specified
const DefaultRadiusKm = 25.0

// degreesToRadians converts degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// HaversineDistance calculates the distance between two points on Earth
// using the Haversine formula. Returns distance in kilometers.
func HaversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return EarthRadiusKm * c
}

// HaversineSQL returns a PostgreSQL expression for calculating distance
// between a point (centerLat, centerLng) and database columns (latCol, lngCol)
func HaversineSQL(latCol, lngCol string, centerLat, centerLng float64) string {
	return fmt.Sprintf(`(
		%f * acos(
			cos(radians(%f)) * cos(radians(%s)) *
			cos(radians(%s) - radians(%f)) +
			sin(radians(%f)) * sin(radians(%s))
		)
	)`, EarthRadiusKm, centerLat, latCol, lngCol, centerLng, centerLat, latCol)
}

// BoundingBox represents a geographic bounding box for pre-filtering
type BoundingBox struct {
	MinLat float64
	MaxLat float64
	MinLng float64
	MaxLng float64
}

// CalculateBoundingBox returns a bounding box around a center point
// that can be used to pre-filter results before applying Haversine
func CalculateBoundingBox(centerLat, centerLng, radiusKm float64) BoundingBox {
	// Angular distance in radians
	angularDistance := radiusKm / EarthRadiusKm

	minLat := centerLat - angularDistance*180/math.Pi
	maxLat := centerLat + angularDistance*180/math.Pi

	// Longitude bounds (accounting for narrowing at higher latitudes)
	latRad := degreesToRadians(centerLat)
	lngDelta := math.Asin(math.Sin(angularDistance) / math.Cos(latRad))
	lngDeltaDeg := lngDelta * 180 / math.Pi

	minLng := centerLng - lngDeltaDeg
	maxLng := centerLng + lngDeltaDeg

	// Handle edge cases
	if minLat < -90 {
		minLat = -90
	}
	if maxLat > 90 {
		maxLat = 90
	}
	if minLng < -180 {
		minLng = -180
	}
	if maxLng > 180 {
		maxLng = 180
	}

	return BoundingBox{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}
}

// BoundingBoxSQL returns a SQL WHERE clause for the bounding box pre-filter
func BoundingBoxSQL(latCol, lngCol string, bbox BoundingBox) string {
	return fmt.Sprintf(`(%s BETWEEN %f AND %f AND %s BETWEEN %f AND %f)`,
		latCol, bbox.MinLat, bbox.MaxLat,
		lngCol, bbox.MinLng, bbox.MaxLng)
}
