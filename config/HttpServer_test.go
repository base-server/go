package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestHttpServer(t *testing.T) {
	var httpServerConfig config.HttpServer

	err := config.Parsing(&httpServerConfig, "./HttpServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if httpServerConfig.SwaggerAddress != "127.0.0.1:10000" {
		t.Errorf("invalid data - SwaggerAddress : (%s)", httpServerConfig.SwaggerAddress)
	}

	if httpServerConfig.SwaggerUri != "/swagger/" {
		t.Errorf("invalid data - SwaggerUri : (%s)", httpServerConfig.SwaggerUri)
	}

	if httpServerConfig.ServerAddress != ":10000" {
		t.Errorf("invalid data - ServerAddress : (%s)", httpServerConfig.ServerAddress)
	}

	if httpServerConfig.ShutdownTimeout != 10 {
		t.Errorf("invalid data - ShutdownTimeout : (%d)", httpServerConfig.ShutdownTimeout)
	}

	if httpServerConfig.Log.Level != "DEBUG" {
		t.Errorf("invalid data - Log.Level : (%s)", httpServerConfig.Log.Level)
	}

	if httpServerConfig.Log.OutputPath != "./log/" {
		t.Errorf("invalid data - Log.OutputPath : (%s)", httpServerConfig.Log.OutputPath)
	}

	if httpServerConfig.Log.FileNamePrefix != "http-server" {
		t.Errorf("invalid data - Log.FileNamePrefix : (%s)", httpServerConfig.Log.FileNamePrefix)
	}

	if httpServerConfig.Log.PrintCallerInfo != true {
		t.Errorf("invalid data - Log.PrintCallerInfo : (%t)", httpServerConfig.Log.PrintCallerInfo)
	}

	if httpServerConfig.Log.ChannelSize != 1024 {
		t.Errorf("invalid data - Log.ChannelSize : (%d)", httpServerConfig.Log.ChannelSize)
	}
}
