package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestSample(t *testing.T) {
	sampleConfig, err := config.Get[config.Sample]("./Sample.config")
	if err != nil {
		t.Fatal(err)
	}

	if sampleConfig.Log.Level != "debug" {
		t.Fatal("invalid -", sampleConfig.Log.Level)
	}

	if sampleConfig.Log.Output != "file" {
		t.Fatal("invalid -", sampleConfig.Log.Output)
	}

	if sampleConfig.Log.File.Name != "./sample" {
		t.Fatal("invalid -", sampleConfig.Log.File.Name)
	}

	if sampleConfig.Log.File.ExtensionName != "log" {
		t.Fatal("invalid -", sampleConfig.Log.File.ExtensionName)
	}

	if sampleConfig.Log.File.AddDate {
		t.Fatal("invalid -", sampleConfig.Log.File.AddDate)
	}

	if sampleConfig.Log.WithCallerInfo == false {
		t.Fatal("invalid -", sampleConfig.Log.WithCallerInfo)
	}
}
