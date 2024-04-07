package main

import (
	"flag"
	"math/rand/v2"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/config"
	"github.com/common-library/go/file"
	"github.com/common-library/go/socket"
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
	configFile := path + "/../config/SocketServer.config"

	socketServerConfig, err := config.Get[config.SocketServer](configFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Remove(socketServerConfig.Log.File.Name + "." + socketServerConfig.Log.File.ExtensionName)

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

	clientJob := func(wg *sync.WaitGroup) {
		defer wg.Done()

		client := socket.Client{}
		defer client.Close()

		if err := client.Connect("tcp", socketServerConfig.Address); err != nil {
			t.Fatal(err)
		}

		if readData, err := client.Read(1024); err != nil {
			t.Fatal(err)
		} else if readData != "greeting" {
			t.Fatalf("invalid data - (%s)", readData)
		}

		writeData := "test-" + strconv.Itoa(rand.IntN(1000))
		if _, err = client.Write(writeData); err != nil {
			t.Fatal(err)
		}

		if readData, err := client.Read(1024); err != nil {
			t.Fatal(err)
		} else if readData != "[response] "+writeData {
			t.Fatalf("invalid data - (%s)", readData)
		}
	}

	wg := sync.WaitGroup{}
	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go clientJob(&wg)
	}
	wg.Wait()

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	for condition.Load() {
		time.Sleep(100 * time.Millisecond)
	}
}
