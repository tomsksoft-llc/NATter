package entity

import (
	"strings"
)

type Route struct {
	Mode     RouteMode      `toml:"MODE" json:"mode"`
	Async    bool           `toml:"ASYNC" json:"async,omitempty"`
	Topic    string         `toml:"TOPIC" json:"topic,omitempty"`
	Endpoint string         `toml:"ENDPOINT" json:"endpoint,omitempty"`
	URI      string         `toml:"URI" json:"uri,omitempty"`
	Batching *RouteBatching `toml:"BATCHING" json:"batching,omitempty"`
}

type RouteBatching struct {
	Timeout  uint32 `toml:"TIMEOUT" json:"timeout"`
	Capacity uint32 `toml:"CAPACITY" json:"capacity"`
}

type RouteMode string

const (
	routeModeSep = "-"
)

func (m RouteMode) Components() *RouteModeComponents {
	comp := make([]string, 3)

	copy(comp, strings.Split(string(m), routeModeSep))

	return &RouteModeComponents{
		Receiver:  comp[0],
		Sender:    comp[1],
		Direction: RouteDirection(comp[2]),
	}
}

type RouteModeComponents struct {
	Receiver  string
	Sender    string
	Direction RouteDirection
}

type RouteDirection string

const (
	RouteDirectionOneway RouteDirection = "oneway"
	RouteDirectionTwoway RouteDirection = "twoway"
)
