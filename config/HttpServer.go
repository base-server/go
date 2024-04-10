package config

import "time"

type HttpServer struct {
	SwaggerAddress string `json:"swagger_address"`
	SwaggerUri     string `json:"swagger_uri"`

	ServerAddress   string        `json:"server_address"`
	ShutdownTimeout time.Duration `json:"shutdownTimeout"`

	Log struct {
		Level  string `json:"level"`
		Output string `json:"output"`
		File   struct {
			Name          string `json:"name"`
			ExtensionName string `json:"extensionName"`
			AddDate       bool   `json:"addDate"`
		} `json:"file"`
		WithCallerInfo bool `json:"withCallerInfo"`
	} `json:"log"`
}
