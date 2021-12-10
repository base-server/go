// Package config provides a struct that can store json type config file
package config

type GrpcServer struct {
	LogLevel          string `json:"log_level"`
	LogOutputPath     string `json:"log_output_path"`
	LogFileNamePrefix string `json:"log_file_name_prefix"`

	Address string `json:"address"`
}

type SocketServer struct {
	LogLevel          string `json:"log_level"`
	LogOutputPath     string `json:"log_output_path"`
	LogFileNamePrefix string `json:"log_file_name_prefix"`

	Address        string `json:"address"`
	ClientPoolSize int    `json:"client_pool_size"`
}
