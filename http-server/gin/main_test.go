package main

import (
	"flag"
	net_http "net/http"
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/common-library/go/file"
	"github.com/common-library/go/http"
)

func TestMain(t *testing.T) {
	const configFile = "../../common/config/config.yaml"

	if err := config.Read(configFile); err != nil {
		t.Fatal(err)
	}

	defer file.Remove(config.Get("http.log.file.name").(string) + "." + config.Get("http.log.file.extensionName").(string))

	condition := atomic.Bool{}
	condition.Store(false)
	go func() {
		os.Args = []string{"test", "-config-file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		condition.Store(true)
		main()
		condition.Store(false)
	}()
	for condition.Load() == false {
		time.Sleep(100 * time.Millisecond)
	}

	port := config.Get("http.gin.address").(string)[1:]

	if response, err := http.Request("http://127.0.0.1:"+port+"/v1/test/id-01?param-1=value-1&param-2=2&param-3=3.3", net_http.MethodGet, map[string][]string{"header-1": {"value-1"}}, "", 3*time.Second, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
	} else if response.Body != `{"id":"id-01","field-1":1,"field-2":"value-2"}` {
		t.Fatalf("invalid Body - (%s)", response.Body)
	}

	if response, err := http.Request("http://127.0.0.1:"+port+"/v1/test", net_http.MethodPost, nil, "", 3*time.Second, "", "", nil); err != nil {
		t.Fatal(err)
	} else if response.StatusCode != net_http.StatusOK {
		t.Fatalf("invalid StatusCode - (%d)", response.StatusCode)
	} else if response.Body != `{"field-1":"value-1"}` {
		t.Fatalf("invalid Body - (%s)", response.Body)
	}

	if response, err := http.Request("http://127.0.0.1:"+port+"/v1/test/id-01", net_http.MethodDelete, nil, "", 3*time.Second, "", "", nil); err != nil {
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
