package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestLongPollingServer(t *testing.T) {
	longPollingServerConfig := config.LongPollingServer{}

	err := config.Parsing(&longPollingServerConfig, "./LongPollingServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if longPollingServerConfig.Address != ":30000" {
		t.Errorf("invalid data - Address : (%s)", longPollingServerConfig.Address)
	}

	if longPollingServerConfig.Timeout != 3600 {
		t.Errorf("invalid data - Timeout : (%d)", longPollingServerConfig.Timeout)
	}

	if longPollingServerConfig.ShutdownTimeout != 10 {
		t.Errorf("invalid data - ShutdownTimeout : (%d)", longPollingServerConfig.ShutdownTimeout)
	}

	if longPollingServerConfig.SubscriptionURI != "/subscription" {
		t.Errorf("invalid data - SubscriptionURI : (%s)", longPollingServerConfig.SubscriptionURI)
	}

	if longPollingServerConfig.PublishURI != "/publish" {
		t.Errorf("invalid data - PublishURI : (%s)", longPollingServerConfig.PublishURI)
	}

	if longPollingServerConfig.FilePersistorInfo.Use != false {
		t.Errorf("invalid data - FilePersistorInfo.Use : (%t)", longPollingServerConfig.FilePersistorInfo.Use)
	}

	if longPollingServerConfig.FilePersistorInfo.FileName != "./file-persistor.txt" {
		t.Errorf("invalid data - FilePersistorInfo.FileName : (%s)", longPollingServerConfig.FilePersistorInfo.FileName)
	}

	if longPollingServerConfig.FilePersistorInfo.WriteBufferSize != 250 {
		t.Errorf("invalid data - FilePersistorInfo.WriteBufferSize : (%d)", longPollingServerConfig.FilePersistorInfo.WriteBufferSize)
	}

	if longPollingServerConfig.FilePersistorInfo.WriteFlushPeriodSeconds != 1 {
		t.Errorf("invalid data - FilePersistorInfo.WriteFlushPeriodSeconds : (%d)", longPollingServerConfig.FilePersistorInfo.WriteFlushPeriodSeconds)
	}

	if longPollingServerConfig.Log.Level != "DEBUG" {
		t.Errorf("invalid data - Log.Level : (%s)", longPollingServerConfig.Log.Level)
	}

	if longPollingServerConfig.Log.OutputPath != "./log/" {
		t.Errorf("invalid data - Log.OutputPath : (%s)", longPollingServerConfig.Log.OutputPath)
	}

	if longPollingServerConfig.Log.FileNamePrefix != "long-polling-server" {
		t.Errorf("invalid data - Log.FileNamePrefix : (%s)", longPollingServerConfig.Log.FileNamePrefix)
	}

	if longPollingServerConfig.Log.PrintCallerInfo != true {
		t.Errorf("invalid data - Log.PrintCallerInfo : (%t)", longPollingServerConfig.Log.PrintCallerInfo)
	}

	if longPollingServerConfig.Log.ChannelSize != 1024 {
		t.Errorf("invalid data - Log.ChannelSize : (%d)", longPollingServerConfig.Log.ChannelSize)
	}
}
