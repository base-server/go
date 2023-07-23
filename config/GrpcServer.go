// Package config provides a struct that can store json type config file
package config

import "github.com/heaven-chp/common-library-go/json"

type GrpcServer struct {
	LogLevel          string `json:"log_level"`
	LogOutputPath     string `json:"log_output_path"`
	LogFileNamePrefix string `json:"log_file_name_prefix"`

	Address string `json:"address"`
}

func (this *GrpcServer) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
