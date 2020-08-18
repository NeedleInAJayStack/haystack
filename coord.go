package haystack

import (
	"fmt"
	"math"
)

// A coordinate is composed of a latitude and longitude
type Coord struct {
	lat float64
	lng float64
}

// Create a new coordinate and check that latitude and longitude conform
func NewCoord(lat float64, lng float64) *Coord {
	lat = math.Min(math.Max(lat, -90.0), 90.0) // Clip to -90 to 90
	lng = math.Min(math.Max(lng, 0.0), 180)    // Clip to 0 to 180
	return &Coord{lat: lat, lng: lng}
}

// Represented as "C(lat,lng)"
func (coord *Coord) toZinc() string {
	result := "C("
	result = result + fmt.Sprintf("%g", coord.lat)
	result = result + ","
	result = result + fmt.Sprintf("%g", coord.lng)
	result = result + ")"

	return result
}
