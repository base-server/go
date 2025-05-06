// Package config provides a struct that can store json type config file
package config

import (
	"github.com/spf13/viper"
)

func Read(file string) error {
	viper.SetConfigFile(file)

	if err := viper.ReadInConfig(); err != nil {
		return err
	} else {
		return nil
	}
}

func Get(key string) any {
	return viper.Get(key)
}
