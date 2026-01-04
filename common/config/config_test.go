package config_test

import (
	"testing"

	"github.com/base-server/go/common/config"
)

func TestRead(t *testing.T) {
	if err := config.Read("./config.yaml"); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	answer := map[string]any{
		"cloudEvents.address":         ":40000",
		"cloudEvents.shutdownTimeout": "10s",

		"gRPC.address":                ":50051",
		"gRPC.log.level":              "debug",
		"gRPC.log.output":             "file",
		"gRPC.log.file.name":          "./grpc-server",
		"gRPC.log.file.extensionName": "log",
		"gRPC.log.file.addDate":       false,
		"gRPC.log.withCallerInfo":     true,

		"http.echo.address":           ":10000",
		"http.gin.address":            ":10001",
		"http.mux.address":            ":10002",
		"http.swaggerAddress":         "127.0.0.1:10000",
		"http.swaggerUri":             "/swagger/",
		"http.shutdownTimeout":        "10s",
		"http.log.level":              "debug",
		"http.log.output":             "file",
		"http.log.file.name":          "./http-server",
		"http.log.file.extensionName": "log",
		"http.log.file.addDate":       false,
		"http.log.withCallerInfo":     true,

		"longPolling.address":                                   ":30000",
		"longPolling.timeoutSeconds":                            3600,
		"longPolling.shutdownTimeout":                           "10s",
		"longPolling.subscriptionURI":                           "/subscription",
		"longPolling.publishURI":                                "/publish",
		"longPolling.filePersistorInfo.use":                     false,
		"longPolling.filePersistorInfo.fileName":                "./file-persistor.txt",
		"longPolling.filePersistorInfo.writeBufferSize":         250,
		"longPolling.filePersistorInfo.writeFlushPeriodSeconds": 1,
		"longPolling.log.level":                                 "debug",
		"longPolling.log.output":                                "file",
		"longPolling.log.file.name":                             "./long-polling-server",
		"longPolling.log.file.extensionName":                    "log",
		"longPolling.log.file.addDate":                          false,
		"longPolling.log.withCallerInfo":                        true,

		"socket.tcp.address":            ":20000",
		"socket.udp.address":            ":20001",
		"socket.clientPoolSize":         1024,
		"socket.log.level":              "debug",
		"socket.log.output":             "file",
		"socket.log.file.name":          "./socket-server",
		"socket.log.file.extensionName": "log",
		"socket.log.file.addDate":       false,
		"socket.log.withCallerInfo":     true,
	}

	if err := config.Read("./config.yaml"); err != nil {
		t.Fatal(err)
	}

	for key, value := range answer {
		if result := config.Get(key); result != value {
			t.Fatal(key, ",", value, ",", result)
		}
	}
}
