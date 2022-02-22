package opensky

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fraktt/locator/internal/app"
)

type client struct {
	url string
}

const allStatesURL = "https://opensky-network.org/api/states/all"

// New creates opensky client
func New() app.Sky {
	return client{
		url: allStatesURL,
	}
}

var _ app.Sky = client{} // check if Client implements app.Sky interface

func (c client) AllPlanes() ([]app.Plane, error) {
	respBody, err := get(c.url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %w", err)
	}

	planes, err := jsonParse(respBody)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	return planes, nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http GET: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("response body read: %w", err)
	}

	return respBody, nil
}

// indexes from https://openskynetwork.github.io/opensky-api/rest.html#response
const (
	icao24Index        = 0
	callSignIndex      = 1
	originCountryIndex = 2
	longitudeIndex     = 5
	latitudeIndex      = 6
)

type Response struct {
	States [][]interface{} `json:"states"`
}

func jsonParse(data []byte) ([]app.Plane, error) {
	var r Response
	if err := json.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	planes := make([]app.Plane, 0, len(r.States))
	for _, state := range r.States {
		maxInd := maxIndex(icao24Index, callSignIndex, originCountryIndex, longitudeIndex, latitudeIndex)
		if maxInd >= len(state) { //to avoid panic on index out of range
			log.Printf("state list is too short to contain all plane properties")
			continue
		}

		icao := state[icao24Index].(string)
		callSign := state[callSignIndex].(string)
		country := state[originCountryIndex].(string)

		long, longOK := state[longitudeIndex].(float64)
		lat, latOK := state[latitudeIndex].(float64)
		if !longOK || !latOK {
			continue // skip planes without coordinates
		}

		planes = append(planes, app.Plane{
			ICAO24:    icao,
			CallSign:  callSign,
			Country:   country,
			Latitude:  lat,
			Longitude: long,
		})
	}

	return planes, nil
}

func maxIndex(indexes ...int) int {
	maxInd := 0
	for _, ind := range indexes {
		if ind > maxInd {
			maxInd = ind
		}
	}
	return maxInd
}
