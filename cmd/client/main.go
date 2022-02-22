package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/fraktt/locator/internal/app"
	"github.com/gorilla/websocket"
)

func main() {
	var (
		addr                        string
		longitude, latitude, radius float64
	)
	flag.StringVar(&addr, "addr", "localhost:8080", "server address")
	flag.Float64Var(&longitude, "long", 0.0, "Longitude")
	flag.Float64Var(&latitude, "lat", 0.0, "Latitude")
	flag.Float64Var(&radius, "rad", 0.0, "Radius")
	flag.Parse()

	log.Println(addr, longitude, latitude, radius)

	// connect to server
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	// send coordinates and radius
	if err := c.WriteJSON(app.ClientMessage{
		Longitude: longitude,
		Latitude:  latitude,
		Radius:    radius,
	}); err != nil {
		log.Println("send json:", err)
		return
	}

	// listen
	var srvMsg app.ServerMessage
	for {
		if err := c.ReadJSON(&srvMsg); err != nil {
			log.Println("receive json:", err)
			return
		}
		fmt.Printf("%d\n", srvMsg.NumberOfPlanes)
	}
}
