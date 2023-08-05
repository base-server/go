package config

import "github.com/heaven-chp/common-library-go/json"

type SocketServer struct {
	Address        string `json:"address"`
	ClientPoolSize int    `json:"client_pool_size"`

	Log struct {
		Level           string `json:"level"`
		OutputPath      string `json:"output_path"`
		FileNamePrefix  string `json:"file_name_prefix"`
		PrintCallerInfo bool   `json:"print_caller_info"`
		ChannelSize     int    `json:"channel_size"`
	} `json:"log"`
}

func (this *SocketServer) parsing(from interface{}) error {
	return json.ToStructFromFile(from.(string), this)
}
