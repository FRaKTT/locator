package app

// ClientMessage contains data client sends to server
type ClientMessage struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Radius    float64 `json:"radius"`
}

// ServerMessage contains data server sends to client
type ServerMessage struct {
	NumberOfPlanes int `json:"number_of_planes"`
}

// Plane properties
type Plane struct {
	ICAO24    string
	CallSign  string
	Country   string
	Longitude float64
	Latitude  float64
}

// Sky - interface for getting planes
type Sky interface {
	AllPlanes() ([]Plane, error)
}

// PlanesCache - interface for caching planes
type PlanesCache interface {
	IsEmpty() bool
	SetPlanes([]Plane)
	Count(func(Plane) bool) int
}
