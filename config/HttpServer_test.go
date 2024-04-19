package config_test

import (
	"testing"

	"github.com/base-server/go/config"
)

func TestHttpServer(t *testing.T) {
	httpServerConfig, err := config.Get[config.HttpServer]("./HttpServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if httpServerConfig.SwaggerAddress != "127.0.0.1:10000" {
		t.Fatal("invalid -", httpServerConfig.SwaggerAddress)
	}

	if httpServerConfig.SwaggerUri != "/swagger/" {
		t.Fatal("invalid -", httpServerConfig.SwaggerUri)
	}

	if httpServerConfig.ServerAddress != ":10000" {
		t.Fatal("invalid -", httpServerConfig.ServerAddress)
	}

	if httpServerConfig.ShutdownTimeout != "10s" {
		t.Fatal("invalid -", httpServerConfig.ShutdownTimeout)
	}

	if httpServerConfig.Log.Level != "debug" {
		t.Fatal("invalid -", httpServerConfig.Log.Level)
	}

	if httpServerConfig.Log.Output != "file" {
		t.Fatal("invalid -", httpServerConfig.Log.Output)
	}

	if httpServerConfig.Log.File.Name != "./http-server" {
		t.Fatal("invalid -", httpServerConfig.Log.File.Name)
	}

	if httpServerConfig.Log.File.ExtensionName != "log" {
		t.Fatal("invalid -", httpServerConfig.Log.File.ExtensionName)
	}

	if httpServerConfig.Log.File.AddDate {
		t.Fatal("invalid -", httpServerConfig.Log.File.AddDate)
	}

	if httpServerConfig.Log.WithCallerInfo == false {
		t.Fatal("invalid -", httpServerConfig.Log.WithCallerInfo)
	}
}
