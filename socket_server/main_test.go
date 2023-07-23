package main

import (
	"flag"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/common-library-go/json"
	"github.com/heaven-chp/common-library-go/socket"
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
	configFile := path + "/../config/socket_server.config"

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
		var client socket.Client
		defer client.Close()

		socketServerConfig := config.SocketServer{}

		err := json.ToStructFromFile(configFile, &socketServerConfig)
		if err != nil {
			t.Fatal(err)
		}

		err = client.Connect("tcp", socketServerConfig.Address)
		if err != nil {
			t.Fatal(err)
		}

		readData, err := client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}
		if readData != "greeting" {
			t.Fatalf("invalid data - (%s)", readData)
		}

		writeData := "test"
		_, err = client.Write(writeData)
		if err != nil {
			t.Fatal(err)
		}

		readData, err = client.Read(1024)
		if err != nil {
			t.Fatal(err)
		}
		if readData != "[response] "+writeData {
			t.Fatalf("invalid data - (%s)", readData)
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
