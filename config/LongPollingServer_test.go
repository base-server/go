package config_test

import (
	"testing"

	"github.com/base-server/go/config"
)

func TestLongPollingServer(t *testing.T) {
	longPollingServerConfig, err := config.Get[config.LongPollingServer]("./LongPollingServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if longPollingServerConfig.Address != ":30000" {
		t.Fatal("invalid -", longPollingServerConfig.Address)
	}

	if longPollingServerConfig.Timeout != 3600 {
		t.Fatal("invalid -", longPollingServerConfig.Timeout)
	}

	if longPollingServerConfig.ShutdownTimeout != 10 {
		t.Fatal("invalid -", longPollingServerConfig.ShutdownTimeout)
	}

	if longPollingServerConfig.SubscriptionURI != "/subscription" {
		t.Fatal("invalid -", longPollingServerConfig.SubscriptionURI)
	}

	if longPollingServerConfig.PublishURI != "/publish" {
		t.Fatal("invalid -", longPollingServerConfig.PublishURI)
	}

	if longPollingServerConfig.FilePersistorInfo.Use != false {
		t.Fatal("invalid -", longPollingServerConfig.FilePersistorInfo.Use)
	}

	if longPollingServerConfig.FilePersistorInfo.FileName != "./file-persistor.txt" {
		t.Fatal("invalid -", longPollingServerConfig.FilePersistorInfo.FileName)
	}

	if longPollingServerConfig.FilePersistorInfo.WriteBufferSize != 250 {
		t.Fatal("invalid -", longPollingServerConfig.FilePersistorInfo.WriteBufferSize)
	}

	if longPollingServerConfig.FilePersistorInfo.WriteFlushPeriodSeconds != 1 {
		t.Fatal("invalid -", longPollingServerConfig.FilePersistorInfo.WriteFlushPeriodSeconds)
	}

	if longPollingServerConfig.Log.Level != "debug" {
		t.Fatal("invalid -", longPollingServerConfig.Log.Level)
	}

	if longPollingServerConfig.Log.Output != "file" {
		t.Fatal("invalid -", longPollingServerConfig.Log.Output)
	}

	if longPollingServerConfig.Log.File.Name != "./long-polling-server" {
		t.Fatal("invalid -", longPollingServerConfig.Log.File.Name)
	}

	if longPollingServerConfig.Log.File.ExtensionName != "log" {
		t.Fatal("invalid -", longPollingServerConfig.Log.File.ExtensionName)
	}

	if longPollingServerConfig.Log.File.AddDate {
		t.Fatal("invalid -", longPollingServerConfig.Log.File.AddDate)
	}

	if longPollingServerConfig.Log.WithCallerInfo == false {
		t.Fatal("invalid -", longPollingServerConfig.Log.WithCallerInfo)
	}
}
