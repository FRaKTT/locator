package app

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	calc "github.com/fraktt/locator/internal/calculations"

	"github.com/gorilla/websocket"
)

type WsClient struct {
	parameters ClientMessage
	conn       *websocket.Conn
	ch         chan int
	id         uint64
}

type App struct {
	sky   Sky
	cache PlanesCache

	clientsMap   map[uint64]*WsClient
	nextClientID uint64
	clientsMx    sync.Mutex
}

const interval = 5 * time.Minute // cache update interval

func New(sky Sky, cache PlanesCache) *App {
	a := &App{
		sky:        sky,
		cache:      cache,
		clientsMap: make(map[uint64]*WsClient),
	}

	go func() { // start to watch the sky and update cache
		// initial cache filling
		log.Println("filling the cache...")
		if err := a.updateCache(); err != nil {
			log.Printf("fill cache: %s\n", err)
		}

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("updating the cache...")
				if err := a.updateCache(); err != nil {
					log.Printf("update cache: %s\n", err)
					break
				}
			}
		}
	}()

	return a
}

// updateCache gets planes, writes them to cache and calculates number of planes for connected clients
func (a *App) updateCache() error {
	planes, err := a.sky.AllPlanes()
	if err != nil {
		return fmt.Errorf("get all planes: %w", err)
	}

	a.cache.SetPlanes(planes)

	a.clientsMx.Lock()
	for _, c := range a.clientsMap {
		c.ch <- a.calcPlanesInRadius(c.parameters)
	}
	a.clientsMx.Unlock()

	return nil
}

func (a *App) calcPlanesInRadius(param ClientMessage) int {
	clientCoords := calc.Coordinates{
		Longitude: param.Longitude,
		Latitude:  param.Latitude,
	}

	return a.cache.Count(func(p Plane) bool {
		planeCoords := calc.Coordinates{
			Longitude: p.Longitude,
			Latitude:  p.Latitude,
		}
		distance := calc.OrthodromicDistance(clientCoords, planeCoords)
		return distance <= param.Radius
	})
}

func (a *App) addClient(msg ClientMessage, conn *websocket.Conn) *WsClient {
	a.clientsMx.Lock()

	id := a.nextClientID
	cl := &WsClient{
		parameters: msg,
		conn:       conn,
		ch:         make(chan int),
		id:         id,
	}
	a.clientsMap[id] = cl

	a.nextClientID++ //todo: number of clients is limited by MaxUint64 = ^uint64(0)

	a.clientsMx.Unlock()
	return cl
}

func (a *App) deleteClient(clientID uint64) {
	a.clientsMx.Lock()
	//todo: close channel and connection?
	delete(a.clientsMap, clientID)
	a.clientsMx.Unlock()
}

var upgrader = websocket.Upgrader{}

func (a *App) GetHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer conn.Close()

		// get client settings
		var clMsg ClientMessage
		if err := conn.ReadJSON(&clMsg); err != nil {
			log.Println("unmarshal:", err)
			return
		}

		// register client
		cl := a.addClient(clMsg, conn)
		defer a.deleteClient(cl.id)
		log.Printf("new client: id %d, parameters: long %v, lat %v, radius %v",
			cl.id, cl.parameters.Longitude, cl.parameters.Latitude, cl.parameters.Radius)

		// send number of planes immediately
		if !a.cache.IsEmpty() {
			n := a.calcPlanesInRadius(clMsg)
			if err := conn.WriteJSON(ServerMessage{NumberOfPlanes: n}); err != nil {
				log.Printf("send json: %s\n", err)
				return
			}
		}

		// periodically send number of planes
		for {
			select {
			case n := <-cl.ch:
				if err := conn.WriteJSON(ServerMessage{NumberOfPlanes: n}); err != nil {
					log.Printf("send json: %s\n", err)
					return
				}
			}
		}
	}
}
