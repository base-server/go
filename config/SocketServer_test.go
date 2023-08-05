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

	if socketServerConfig.Address != ":20000" {
		t.Errorf("invalid data - Address : (%s)", socketServerConfig.Address)
	}

	if socketServerConfig.ClientPoolSize != 1024 {
		t.Errorf("invalid data - ClientPoolSize : (%d)", socketServerConfig.ClientPoolSize)
	}

	if socketServerConfig.Log.Level != "DEBUG" {
		t.Errorf("invalid data - Log.Level : (%s)", socketServerConfig.Log.Level)
	}

	if socketServerConfig.Log.OutputPath != "./log/" {
		t.Errorf("invalid data - Log.OutputPath : (%s)", socketServerConfig.Log.OutputPath)
	}

	if socketServerConfig.Log.FileNamePrefix != "socket-server" {
		t.Errorf("invalid data - Log.FileNamePrefix : (%s)", socketServerConfig.Log.FileNamePrefix)
	}

	if socketServerConfig.Log.PrintCallerInfo != true {
		t.Errorf("invalid data - Log.PrintCallerInfo : (%t)", socketServerConfig.Log.PrintCallerInfo)
	}

	if socketServerConfig.Log.ChannelSize != 1024 {
		t.Errorf("invalid data - Log.ChannelSize : (%d)", socketServerConfig.Log.ChannelSize)
	}
}
