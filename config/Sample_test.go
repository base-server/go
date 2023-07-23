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

	if sampleConfig.Field1 != 1 {
		t.Errorf("invalid field1 value - (%d)", sampleConfig.Field1)
	}

	if sampleConfig.Field2 != "value2" {
		t.Errorf("invalid field2 value - (%s)", sampleConfig.Field2)
	}
}
