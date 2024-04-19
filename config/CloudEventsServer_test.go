package config_test

import (
	"testing"

	"github.com/base-server/go/config"
)

func TestCloudEventsServer(t *testing.T) {
	if cloudEventsServerConfig, err := config.Get[config.CloudEventsServer]("./CloudEventsServer.config"); err != nil {
		t.Fatal(err)
	} else if cloudEventsServerConfig.Address != ":40000" {
		t.Fatal("invalid -", cloudEventsServerConfig.Address)
	} else if cloudEventsServerConfig.ShutdownTimeout != "10s" {
		t.Fatal("invalid -", cloudEventsServerConfig.ShutdownTimeout)
	}
}
