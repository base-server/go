package main

import (
	"context"
	"flag"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/config"
	"github.com/common-library/go/file"
	"github.com/common-library/go/grpc"
	"github.com/common-library/go/grpc/sample"
)

func TestMain1(t *testing.T) {
	os.Args = []string{"test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if err := (&Main{}).Run(); err.Error() != "invalid flag" {
		t.Fatal(err)
	}
}

func TestMain2(t *testing.T) {
	os.Args = []string{"test", "-config_file=invalid"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if err := (&Main{}).Run(); err.Error() != "open invalid: no such file or directory" {
		t.Fatal(err)
	}
}

func TestMain3(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	configFile := path + "/../config/GrpcServer.config"

	grpcServerConfig, err := config.Get[config.GrpcServer](configFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Remove(grpcServerConfig.Log.File.Name + "." + grpcServerConfig.Log.File.ExtensionName)

	condition := atomic.Bool{}
	condition.Store(true)
	go func() {
		defer condition.Store(false)

		os.Args = []string{"test", "-config_file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()
	}()
	time.Sleep(200 * time.Millisecond)

	func() {
		connection, err := grpc.GetConnection(grpcServerConfig.Address)
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
}
