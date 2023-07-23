package main

import (
	"context"
	"flag"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/common-library-go/grpc"
	"github.com/heaven-chp/common-library-go/grpc/sample"
	"github.com/heaven-chp/common-library-go/json"
)

func TestMain1(t *testing.T) {
	os.Args = []string{"test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main := Main{}
	err := main.Run()
	if err.Error() != "invalid flag" {
		t.Error(err)
	}
}

func TestMain2(t *testing.T) {
	os.Args = []string{"test", "-config_file=invalid"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main := Main{}
	err := main.Run()
	if err.Error() != "open invalid: no such file or directory" {
		t.Error(err)
	}
}

func TestMain3(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	configFile := path + "/../config/grpc_server.config"

	sleep := atomic.Bool{}
	sleep.Store(true)
	condition := atomic.Bool{}
	condition.Store(false)
	go func() {
		os.Args = []string{"test", "-config_file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		sleep.Store(false)
		condition.Store(true)
		main()
		condition.Store(false)
	}()
	for sleep.Load() {
		time.Sleep(100 * time.Millisecond)
	}

	{
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(5)*time.Second)
		defer cancel()

		grpcServerConfig := config.GrpcServer{}

		err := json.ToStructFromFile(configFile, &grpcServerConfig)
		if err != nil {
			t.Fatal(err)
		}

		connection, err := grpc.GetConnection(grpcServerConfig.Address)
		if err != nil {
			t.Fatal(err)
		}
		defer connection.Close()

		client := sample.NewSampleClient(connection)

		request := sample.Request{Data1: 1, Data2: "abc"}
		reply, err := client.Func(ctx, &request)
		if err != nil {
			t.Fatal(err)
		}

		if reply.Data1 != 1 || reply.Data2 != "abc" {
			t.Fatalf("invalid reply - (%d)(%s)", reply.Data1, reply.Data2)
		}
	}

	err = syscall.Kill(os.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Error(err)
	}
	for condition.Load() {
		time.Sleep(100 * time.Millisecond)
	}
}
