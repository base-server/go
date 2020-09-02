package config

import (
	"github.com/heaven-chp/common-library-go/json"
	"testing"
)

func TestGrpcServer(t *testing.T) {
	var grpcServer GrpcServer

	err := json.ToStructFromFile("./grpc_server.config", &grpcServer)
	if err != nil {
		t.Error(err)
	}

	if grpcServer.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", grpcServer.LogLevel)
	}

	if grpcServer.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", grpcServer.LogOutputPath)
	}

	if grpcServer.LogFileNamePrefix != "grpc_server" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", grpcServer.LogFileNamePrefix)
	}

	if grpcServer.Address != "127.0.0.1:50051" {
		t.Errorf("invalid data - Address : (%s)", grpcServer.Address)
	}
}

func TestSocketServer(t *testing.T) {
	var socketServer SocketServer

	err := json.ToStructFromFile("./socket_server.config", &socketServer)
	if err != nil {
		t.Error(err)
	}

	if socketServer.LogLevel != "DEBUG" {
		t.Errorf("invalid data - LogLevel : (%s)", socketServer.LogLevel)
	}

	if socketServer.LogOutputPath != "./log/" {
		t.Errorf("invalid data - LogOutputPath : (%s)", socketServer.LogOutputPath)
	}

	if socketServer.LogFileNamePrefix != "socket_server" {
		t.Errorf("invalid data - LogFileNamePrefix : (%s)", socketServer.LogFileNamePrefix)
	}

	if socketServer.Address != "127.0.0.1:11111" {
		t.Errorf("invalid data - Address : (%s)", socketServer.Address)
	}

	if socketServer.ClientPoolSize != 1024 {
		t.Errorf("invalid data - ClientPoolSize : (%d)", socketServer.ClientPoolSize)
	}
}
