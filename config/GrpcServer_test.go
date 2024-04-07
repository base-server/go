package config_test

import (
	"testing"

	"github.com/base-server/go/config"
)

func TestGrpcServer(t *testing.T) {
	grpcServerConfig, err := config.Get[config.GrpcServer]("./GrpcServer.config")
	if err != nil {
		t.Fatal(err)
	}

	if grpcServerConfig.Address != ":50051" {
		t.Fatal("invalid -", grpcServerConfig.Address)
	}

	if grpcServerConfig.Log.Level != "debug" {
		t.Fatal("invalid -", grpcServerConfig.Log.Level)
	}

	if grpcServerConfig.Log.Output != "file" {
		t.Fatal("invalid -", grpcServerConfig.Log.Output)
	}

	if grpcServerConfig.Log.File.Name != "./grpc-server" {
		t.Fatal("invalid -", grpcServerConfig.Log.File.Name)
	}

	if grpcServerConfig.Log.File.ExtensionName != "log" {
		t.Fatal("invalid -", grpcServerConfig.Log.File.ExtensionName)
	}

	if grpcServerConfig.Log.File.AddDate {
		t.Fatal("invalid -", grpcServerConfig.Log.File.AddDate)
	}

	if grpcServerConfig.Log.WithCallerInfo == false {
		t.Fatal("invalid -", grpcServerConfig.Log.WithCallerInfo)
	}
}
