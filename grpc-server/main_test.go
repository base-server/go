package main

import (
	"context"
	"flag"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/common-library/go/file"
	"github.com/common-library/go/grpc"
	"github.com/common-library/go/grpc/sample"
)

func TestMain(t *testing.T) {
	const configFile = "../common/config/config.yaml"

	// Load config file first
	if err := config.Read(configFile); err != nil {
		t.Fatal(err)
	}

	condition := atomic.Bool{}
	condition.Store(true)
	go func() {
		defer condition.Store(false)

		os.Args = []string{"test", "-config-file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()
	}()
	time.Sleep(500 * time.Millisecond)

	func() {
		connection, err := grpc.GetConnection(config.Get("gRPC.address").(string))
		if err != nil {
			t.Fatal(err)
		}
		defer connection.Close()

		client := sample.NewSampleClient(connection)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		if reply, err := client.Func1(ctx, &sample.Request{Data1: 1, Data2: "abc"}); err != nil {
			t.Fatal(err)
		} else if reply.Data1 != 1 || reply.Data2 != "abc" {
			t.Fatalf("invalid reply - (%d)(%s)", reply.Data1, reply.Data2)
		}
	}()

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	for condition.Load() {
		time.Sleep(100 * time.Millisecond)
	}

	file.Remove(config.Get("gRPC.log.file.name").(string) + "." + config.Get("gRPC.log.file.extensionName").(string))
}
