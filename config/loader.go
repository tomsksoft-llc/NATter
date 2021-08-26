package config

import (
	"github.com/spf13/viper"
)

func Load(env string) error {
	viper.SetConfigFile(env)

	return viper.ReadInConfig()
}

func String(name string) string {
	return viper.GetString(name)
}

func StringSlice(name string) []string {
	return viper.GetStringSlice(name)
}

func Bool(name string) bool {
	return viper.GetBool(name)
}
