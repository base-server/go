package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestSocketServer(t *testing.T) {
	socketServerConfig := config.SocketServer{}

	err := config.Parsing(&socketServerConfig, "./SocketServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if socketServerConfig.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", socketServerConfig.LogLevel)
	}

	if socketServerConfig.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", socketServerConfig.LogOutputPath)
	}

	if socketServerConfig.LogFileNamePrefix != "socket_server" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", socketServerConfig.LogFileNamePrefix)
	}

	if socketServerConfig.Address != ":20000" {
		t.Errorf("invalid data - Address : (%s)", socketServerConfig.Address)
	}

	if socketServerConfig.ClientPoolSize != 1024 {
		t.Errorf("invalid data - ClientPoolSize : (%d)", socketServerConfig.ClientPoolSize)
	}
}
