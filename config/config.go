// Package config provides a struct that can store json type config file
package config

import (
	"github.com/common-library/go/json"
)

func Get[T any](fileName string) (T, error) {
	return json.ConvertFromFile[T](fileName)
}
