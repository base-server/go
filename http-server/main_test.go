package main

import (
	"flag"
	net_http "net/http"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/heaven-chp/common-library-go/http"
)

func TestMain1(t *testing.T) {
	os.Args = []string{"test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main := Main{}
	err := main.Run()
	if err == nil || err.Error() != "invalid flag" {
		t.Errorf("invalid error - (%#v)", err)
	}
}

func TestMain2(t *testing.T) {
	os.Args = []string{"test", "-config_file=invalid"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	main := Main{}
	err := main.Run()
	if err.Error() != "open invalid: no such file or directory" {
		t.Errorf("invalid error - (%s)", err.Error())
	}
}

func TestMain3(t *testing.T) {
	var condition atomic.Bool
	condition.Store(false)

	go func() {
		path, err := os.Getwd()
		if err != nil {
			t.Error(err)
		}

		os.Args = []string{"test", "-config_file=" + path + "/../config/HttpServer.config"}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		condition.Store(true)
		main()
		condition.Store(false)
	}()
	for condition.Load() == false {
		time.Sleep(100 * time.Millisecond)
	}

	{
		response, err := http.Request("http://127.0.0.1:10000/v1/test/id-01?param-1=value-1&param-2=2&param-3=3.3", net_http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 3, "", "")
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != net_http.StatusOK {
			t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
		}

		if response.Body != `{"id":"id-01","field-1":1,"field-2":"value-2"}` {
			t.Fatalf("invalid Body - (%s)", response.Body)
		}
	}

	{
		response, err := http.Request("http://127.0.0.1:10000/v1/test", net_http.MethodPost, nil, "", 3, "", "")
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != net_http.StatusOK {
			t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
		}

		if response.Body != `{"field-1":"value-1"}` {
			t.Fatalf("invalid Body - (%s)", response.Body)
		}
	}

	{
		response, err := http.Request("http://127.0.0.1:10000/v1/test/id-01", net_http.MethodDelete, nil, "", 3, "", "")
		if err != nil {
			t.Fatal(err)
		}

		if response.StatusCode != net_http.StatusNoContent {
			t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
		}
	}

	err := syscall.Kill(os.Getpid(), syscall.SIGTERM)
	if err != nil {
		t.Error(err)
	}
	for condition.Load() {
		time.Sleep(100 * time.Millisecond)
	}
}
