package main

import (
	"flag"
	net_http "net/http"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/config"
	"github.com/common-library/go/file"
	"github.com/common-library/go/http"
)

func TestMain1(t *testing.T) {
	os.Args = []string{"test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if err := (&Main{}).Run(); err.Error() != "invalid flag" {
		t.Errorf("invalid error - (%#v)", err)
	}
}

func TestMain2(t *testing.T) {
	os.Args = []string{"test", "-config_file=invalid"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if err := (&Main{}).Run(); err.Error() != "open invalid: no such file or directory" {
		t.Errorf("invalid error - (%s)", err.Error())
	}
}

func TestMain3(t *testing.T) {
	var condition atomic.Bool
	condition.Store(false)

	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	configFile := path + "/../config/HttpServer.config"

	if httpServerConfig, err := config.Get[config.HttpServer](configFile); err != nil {
		t.Fatal(err)
	} else {
		defer file.Remove(httpServerConfig.Log.File.Name + "." + httpServerConfig.Log.File.ExtensionName)
	}

	go func() {
		os.Args = []string{"test", "-config_file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		condition.Store(true)
		main()
		condition.Store(false)
	}()
	for condition.Load() == false {
		time.Sleep(100 * time.Millisecond)
	}

	if response, err := http.Request("http://127.0.0.1:10000/v1/test/id-01?param-1=value-1&param-2=2&param-3=3.3", net_http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 3, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
	} else if response.Body != `{"id":"id-01","field-1":1,"field-2":"value-2"}` {
		t.Fatalf("invalid Body - (%s)", response.Body)
	}

	if response, err := http.Request("http://127.0.0.1:10000/v1/test", net_http.MethodPost, nil, "", 3, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
	} else if response.Body != `{"field-1":"value-1"}` {
		t.Fatalf("invalid Body - (%s)", response.Body)
	}

	if response, err := http.Request("http://127.0.0.1:10000/v1/test/id-01", net_http.MethodDelete, nil, "", 3, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != net_http.StatusNoContent {
		t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
	}

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Error(err)
	}

	for condition.Load() {
		time.Sleep(100 * time.Millisecond)
	}
}
