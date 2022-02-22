package cache

import (
	"sync"

	"github.com/fraktt/locator/internal/app"
)

type cache struct {
	planes []app.Plane
	mx     sync.Mutex
}

// New creates cache
func New() *cache {
	return &cache{}
}

var _ app.PlanesCache = &cache{} // check if cache implements app.PlanesCache interface

func (c *cache) IsEmpty() bool {
	return c.planes == nil
}

func (c *cache) SetPlanes(planes []app.Plane) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.planes = make([]app.Plane, len(planes))
	copy(c.planes, planes)
}

func (c *cache) GetPlanes() []app.Plane {
	c.mx.Lock()
	defer c.mx.Unlock()

	planes := make([]app.Plane, len(c.planes))
	copy(planes, c.planes)
	return planes
}

func (c *cache) Count(predicat func(app.Plane) bool) int {
	c.mx.Lock()
	defer c.mx.Unlock()

	var n int
	for _, p := range c.planes {
		if predicat(p) {
			n++
		}
	}
	return n
}
