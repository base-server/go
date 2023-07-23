package config_test

import (
	"testing"

	"github.com/heaven-chp/base-server-go/config"
)

func TestGrpcServer(t *testing.T) {
	var grpcServerConfig config.GrpcServer

	err := config.Parsing(&grpcServerConfig, "./GrpcServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if grpcServerConfig.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", grpcServerConfig.LogLevel)
	}

	if grpcServerConfig.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", grpcServerConfig.LogOutputPath)
	}

	if grpcServerConfig.LogFileNamePrefix != "grpc_server" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", grpcServerConfig.LogFileNamePrefix)
	}

	if grpcServerConfig.Address != ":50051" {
		t.Errorf("invalid data - Address : (%s)", grpcServerConfig.Address)
	}
}
