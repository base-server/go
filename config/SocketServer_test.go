package config_test

import (
	"testing"

	"github.com/base-server/go/config"
)

func TestSocketServer(t *testing.T) {
	socketServerConfig, err := config.Get[config.SocketServer]("./SocketServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if socketServerConfig.Address != ":20000" {
		t.Fatal("invalid -", socketServerConfig.Address)
	}

	if socketServerConfig.ClientPoolSize != 1024 {
		t.Fatal("invalid -", socketServerConfig.ClientPoolSize)
	}

	if socketServerConfig.Log.Level != "debug" {
		t.Fatal("invalid -", socketServerConfig.Log.Level)
	}

	if socketServerConfig.Log.Output != "file" {
		t.Fatal("invalid -", socketServerConfig.Log.Output)
	}

	if socketServerConfig.Log.File.Name != "./socket-server" {
		t.Fatal("invalid -", socketServerConfig.Log.File.Name)
	}

	if socketServerConfig.Log.File.ExtensionName != "log" {
		t.Fatal("invalid -", socketServerConfig.Log.File.ExtensionName)
	}

	if socketServerConfig.Log.File.AddDate {
		t.Fatal("invalid -", socketServerConfig.Log.File.AddDate)
	}

	if socketServerConfig.Log.WithCallerInfo == false {
		t.Fatal("invalid -", socketServerConfig.Log.WithCallerInfo)
	}
}
