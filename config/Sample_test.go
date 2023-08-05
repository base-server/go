package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestSample(t *testing.T) {
	sampleConfig := config.Sample{}

	err := config.Parsing(&sampleConfig, "./Sample.config")
	if err != nil {
		t.Fatal(err)
	}

	if sampleConfig.Log.Level != "DEBUG" {
		t.Errorf("invalid data - Log.Level : (%s)", sampleConfig.Log.Level)
	}

	if sampleConfig.Log.OutputPath != "./log/" {
		t.Errorf("invalid data - Log.OutputPath : (%s)", sampleConfig.Log.OutputPath)
	}

	if sampleConfig.Log.FileNamePrefix != "sample" {
		t.Errorf("invalid data - Log.FileNamePrefix : (%s)", sampleConfig.Log.FileNamePrefix)
	}

	if sampleConfig.Log.PrintCallerInfo != true {
		t.Errorf("invalid data - Log.PrintCallerInfo : (%t)", sampleConfig.Log.PrintCallerInfo)
	}

	if sampleConfig.Log.ChannelSize != 1024 {
		t.Errorf("invalid data - Log.ChannelSize : (%d)", sampleConfig.Log.ChannelSize)
	}
}
