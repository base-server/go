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

	if httpServerConfig.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", httpServerConfig.LogLevel)
	}

	if httpServerConfig.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", httpServerConfig.LogOutputPath)
	}

	if httpServerConfig.LogFileNamePrefix != "http_server" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", httpServerConfig.LogFileNamePrefix)
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
}
