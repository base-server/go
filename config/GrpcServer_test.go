package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestGrpcServer(t *testing.T) {
	grpcServerConfig := config.GrpcServer{}

	err := config.Parsing(&grpcServerConfig, "./GrpcServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if grpcServerConfig.Address != ":50051" {
		t.Errorf("invalid data - Address : (%s)", grpcServerConfig.Address)
	}

	if grpcServerConfig.Log.Level != "DEBUG" {
		t.Errorf("invalid data - Log.Level : (%s)", grpcServerConfig.Log.Level)
	}

	if grpcServerConfig.Log.OutputPath != "./log/" {
		t.Errorf("invalid data - Log.OutputPath : (%s)", grpcServerConfig.Log.OutputPath)
	}

	if grpcServerConfig.Log.FileNamePrefix != "grpc-server" {
		t.Errorf("invalid data - Log.FileNamePrefix : (%s)", grpcServerConfig.Log.FileNamePrefix)
	}

	if grpcServerConfig.Log.PrintCallerInfo != true {
		t.Errorf("invalid data - Log.PrintCallerInfo : (%t)", grpcServerConfig.Log.PrintCallerInfo)
	}

	if grpcServerConfig.Log.ChannelSize != 1024 {
		t.Errorf("invalid data - Log.ChannelSize : (%d)", grpcServerConfig.Log.ChannelSize)
	}
}
