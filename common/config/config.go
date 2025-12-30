// Package config provides a struct that can store json type config file
package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	once    sync.Once
	readErr error
)

func Read(file string) error {
	once.Do(func() {
		viper.SetConfigFile(file)
		readErr = viper.ReadInConfig()
	})

	return readErr
}

func Get(key string) any {
	return viper.Get(key)
}

// Reset resets the config state. This is only for testing purposes.
func Reset() {
	once = sync.Once{}
	readErr = nil
	viper.Reset()
}
