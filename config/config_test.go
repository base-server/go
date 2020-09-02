package config

import (
	"github.com/heaven-chp/common-library-go/json"
	"testing"
)

func TestSocketServer(t *testing.T) {
	var socketServer SocketServer

	err := json.ToStructFromFile("./socket_server.config", &socketServer)
	if err != nil {
		t.Error(err)
	}

	if socketServer.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", socketServer.LogLevel)
	}

	if socketServer.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", socketServer.LogOutputPath)
	}

	if socketServer.LogFileNamePrefix != "socket_server" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", socketServer.LogFileNamePrefix)
	}

	if socketServer.Address != "127.0.0.1:11111" {
		t.Errorf("invalid data - Address : (%s)", socketServer.Address)
	}

	if socketServer.ClientPoolSize != 1024 {
		t.Errorf("invalid data - ClientPoolSize : (%d)", socketServer.ClientPoolSize)
	}
}
