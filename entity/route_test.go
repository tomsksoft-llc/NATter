package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteModeComponents(t *testing.T) {
	mode := RouteMode("receiver-sender-oneway")

	comp := mode.Components()

	assert.Equal(t, &RouteModeComponents{
		Receiver:  "receiver",
		Sender:    "sender",
		Direction: RouteDirectionOneway,
	}, comp)

	mode = RouteMode("receiver-sender-direction")

	comp = mode.Components()

	assert.Equal(t, &RouteModeComponents{
		Receiver:  "receiver",
		Sender:    "sender",
		Direction: RouteDirection("direction"),
	}, comp)

	mode = RouteMode("receiver-sender")

	comp = mode.Components()

	assert.Equal(t, &RouteModeComponents{
		Receiver:  "receiver",
		Sender:    "sender",
		Direction: RouteDirection(""),
	}, comp)

	mode = RouteMode("receiver-sender-twoway-trash-trash")

	comp = mode.Components()

	assert.Equal(t, &RouteModeComponents{
		Receiver:  "receiver",
		Sender:    "sender",
		Direction: RouteDirectionTwoway,
	}, comp)
}
