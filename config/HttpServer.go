package config

import "github.com/heaven-chp/common-library-go/json"

type HttpServer struct {
	LogLevel          string `json:"log_level"`
	LogOutputPath     string `json:"log_output_path"`
	LogFileNamePrefix string `json:"log_file_name_prefix"`

	SwaggerAddress string `json:"swagger_address"`
	SwaggerUri     string `json:"swagger_uri"`

	ServerAddress   string `json:"server_address"`
	ShutdownTimeout uint64 `json:"shutdownTimeout"`
}

func (this *HttpServer) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
