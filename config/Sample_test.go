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

	if sampleConfig.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", sampleConfig.LogLevel)
	}

	if sampleConfig.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", sampleConfig.LogOutputPath)
	}

	if sampleConfig.LogFileNamePrefix != "sample" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", sampleConfig.LogFileNamePrefix)
	}
}
