package config

import "github.com/heaven-chp/common-library-go/json"

type SocketServer struct {
	LogLevel          string `json:"log_level"`
	LogOutputPath     string `json:"log_output_path"`
	LogFileNamePrefix string `json:"log_file_name_prefix"`

	Address        string `json:"address"`
	ClientPoolSize int    `json:"client_pool_size"`
}

func (this *SocketServer) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
