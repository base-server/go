// Package config provides a struct that can store json type config file
package config

type common interface {
	parsing(from interface{}) error
}

func Parsing(to common, from interface{}) error {
	return to.parsing(from)
}
