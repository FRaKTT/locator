package calculations_test

import (
	"math"
	"testing"

	"github.com/fraktt/locator/internal/calculations"
	"github.com/stretchr/testify/assert"
)

const (
	earthR           = 6400  // 6400 kilometers
	defaultTolerance = 0.001 // 1 meter
)

type testCase struct {
	name             string
	c1, c2           calculations.Coordinates
	expectedDistance float64
	tolerance        float64
}

func TestOrthodromicDistance(t *testing.T) {
	cases := [...]testCase{
		{
			name:             "poles distance",
			c1:               calculations.Coordinates{Longitude: 0, Latitude: 90},
			c2:               calculations.Coordinates{Longitude: 0, Latitude: -90},
			expectedDistance: earthR * math.Pi, // half the circumference
		},
		{
			name:             "poles distance, not affecting longitude",
			c1:               calculations.Coordinates{Longitude: 78, Latitude: 90},
			c2:               calculations.Coordinates{Longitude: 32, Latitude: -90},
			expectedDistance: earthR * math.Pi, // half the circumference
		},
		{
			name:             "diametrically opposite on equator",
			c1:               calculations.Coordinates{Longitude: 0, Latitude: 0},
			c2:               calculations.Coordinates{Longitude: 180, Latitude: 0},
			expectedDistance: earthR * math.Pi, // half the circumference
		},

		{
			name:             "1 degree on equator",
			c1:               calculations.Coordinates{Longitude: 0, Latitude: 0},
			c2:               calculations.Coordinates{Longitude: 1, Latitude: 0},
			expectedDistance: earthR * math.Pi / 180, // 1 degree
		},
		{
			name:             "2 degrees on equator",
			c1:               calculations.Coordinates{Longitude: -179, Latitude: 0},
			c2:               calculations.Coordinates{Longitude: 179, Latitude: 0},
			expectedDistance: earthR * math.Pi / 180 * 2, // 2 degrees
		},

		{
			name:             "near the pole, 2 degrees diametrically opposite",
			c1:               calculations.Coordinates{Longitude: -90, Latitude: 89},
			c2:               calculations.Coordinates{Longitude: 90, Latitude: 89},
			expectedDistance: earthR * math.Pi / 180 * 2, // 2 degrees
		},

		{
			name:             "from Moscow to Peter",
			c1:               calculations.Coordinates{Longitude: 37.620430, Latitude: 55.754052},
			c2:               calculations.Coordinates{Longitude: 30.315365, Latitude: 59.938960},
			expectedDistance: 636, // according to maps.yandex.ru
			tolerance:        2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tolerance := defaultTolerance
			if tc.tolerance != 0 {
				tolerance = tc.tolerance
			}

			distance := calculations.OrthodromicDistance(tc.c1, tc.c2)
			delta := math.Abs(distance - tc.expectedDistance)
			assert.Less(t, delta, tolerance)
		})
	}
}
