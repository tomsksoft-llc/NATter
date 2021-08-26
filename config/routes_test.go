package config

import (
	"strings"
	"testing"

	"NATter/entity"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	viper.SetConfigType("toml")

	err := viper.ReadConfig(strings.NewReader(`
		[[ROUTES]]
		MODE="broker-http-twoway"
		TOPIC="topic1"
		ENDPOINT="http://localhost:8081/path1"
		[ROUTES.BATCHING]
		TIMEOUT=30
		CAPACITY=5
		
		[[ROUTES]]
		MODE="http-broker-twoway"
		ASYNC=false
		TOPIC="topic2"
		URI="/path2"
	`))

	assert.Nil(t, err)

	res, err := Routes()

	assert.Nil(t, err)
	assert.Equal(t, []*entity.Route{
		{
			Mode:     entity.RouteMode("broker-http-twoway"),
			Topic:    "topic1",
			Endpoint: "http://localhost:8081/path1",
			Batching: &entity.RouteBatching{
				Timeout:  30,
				Capacity: 5,
			},
		},
		{
			Mode:  entity.RouteMode("http-broker-twoway"),
			Async: false,
			Topic: "topic2",
			URI:   "/path2",
		},
	}, res)
}
