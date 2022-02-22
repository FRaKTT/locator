package calculations

import (
	"math"
)

// Coordinates in degrees
type Coordinates struct {
	Longitude float64
	Latitude  float64
}

// SphericalCoordinates in radians
type SphericalCoordinates struct {
	phi   float64
	theta float64
}

const EarthR = 6400 // approximate radius of a spherical earth (in a vacuum;))

// coordsToSpherical transforms coordinates (longitude and latitude) in degrees
// to spherical coordinates in radians
func coordsToSpherical(c Coordinates) SphericalCoordinates {
	return SphericalCoordinates{
		phi:   c.Longitude * math.Pi / 180,
		theta: (90 - c.Latitude) * math.Pi / 180,
	}
}

// angle between two directions, in radians
func angle(s1, s2 SphericalCoordinates) float64 {
	return math.Acos(
		math.Sin(s1.theta)*math.Sin(s2.theta)*math.Cos(s1.phi-s2.phi) +
			math.Cos(s1.theta)*math.Cos(s2.theta),
	)
}

// OrthodromicDistance returns distance in kilometers between two coordinates
func OrthodromicDistance(c1, c2 Coordinates) float64 {
	s1, s2 := coordsToSpherical(c1), coordsToSpherical(c2)
	return angle(s1, s2) * EarthR
}
