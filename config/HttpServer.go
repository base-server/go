package config

import "github.com/heaven-chp/common-library-go/json"

type HttpServer struct {
	SwaggerAddress string `json:"swagger_address"`
	SwaggerUri     string `json:"swagger_uri"`

	ServerAddress   string `json:"server_address"`
	ShutdownTimeout uint64 `json:"shutdownTimeout"`

	Log struct {
		Level           string `json:"level"`
		OutputPath      string `json:"output_path"`
		FileNamePrefix  string `json:"file_name_prefix"`
		PrintCallerInfo bool   `json:"print_caller_info"`
		ChannelSize     int    `json:"channel_size"`
	} `json:"log"`
}

func (this *HttpServer) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
