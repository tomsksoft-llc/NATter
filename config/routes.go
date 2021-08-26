package config

import (
	"NATter/entity"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func Routes() (routes []*entity.Route, err error) {
	err = viper.UnmarshalKey("ROUTES", &routes, func(ms *mapstructure.DecoderConfig) {
		ms.TagName = "toml"
	})

	if err != nil {
		return nil, err
	}

	return routes, nil
}
