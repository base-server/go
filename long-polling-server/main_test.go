package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/config"
	"github.com/common-library/go/file"
	long_polling "github.com/common-library/go/long-polling"
)

func subscription(t *testing.T, configFile string, request long_polling.SubscriptionRequest, count int, data string) (int64, string) {
	longPollingServerConfig, err := config.Get[config.LongPollingServer](configFile)
	if err != nil {
		t.Fatal(err)
	}

	response, err := long_polling.Subscription("http://"+longPollingServerConfig.Address+"/subscription", nil, request, "", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid status code - (%d)(%s)", response.StatusCode, http.StatusText(response.StatusCode))
	}

	if len(response.Events) != count {
		t.Fatalf("invalid count - (%d)(%d)", len(response.Events), count)
	}

	for _, event := range response.Events {
		if event.Category != request.Category {
			t.Fatalf("invalid category - (%s)(%s)", event.Category, request.Category)
		}

		if event.Data != data {
			t.Fatalf("invalid data - (%s)", event.Data)
		}
	}

	return response.Events[len(response.Events)-1].Timestamp, response.Events[len(response.Events)-1].ID
}

func publish(t *testing.T, configFile, category, data string) {
	longPollingServerConfig, err := config.Get[config.LongPollingServer](configFile)
	if err != nil {
		t.Fatal(err)
	}

	request := long_polling.PublishRequest{Category: category, Data: data}
	response, err := long_polling.Publish("http://"+longPollingServerConfig.Address+longPollingServerConfig.PublishURI, 10, nil, request, "", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("invalid status code - (%d)(%s)", response.StatusCode, http.StatusText(response.StatusCode))
	}

	if response.Body != `{"success": true}` {
		t.Fatalf("invalid body- (%s)", response.Body)
	}
}

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
	configFile := path + "/../config/LongPollingServer.config"

	if longPollingServerConfig, err := config.Get[config.LongPollingServer](configFile); err != nil {
		t.Fatal(err)
	} else {
		defer file.Remove(longPollingServerConfig.Log.File.Name + "." + longPollingServerConfig.Log.File.ExtensionName)
	}

	condition := atomic.Bool{}
	condition.Store(false)
	go func() {
		os.Args = []string{"test", "-config_file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		condition.Store(true)
		main()
		condition.Store(false)
	}()
	time.Sleep(200 * time.Millisecond)

	wg := new(sync.WaitGroup)

	clientJob := func(category, data string) {
		println(category, data)
		defer wg.Done()

		publish(t, configFile, category, data)
		timestamp, id := subscription(t, configFile, long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: 1}, 1, data)

		publish(t, configFile, category, data)
		publish(t, configFile, category, data)
		subscription(t, configFile, long_polling.SubscriptionRequest{Category: category, Timeout: 300, SinceTime: timestamp, LastID: id}, 2, data)
	}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go clientJob("category-"+strconv.Itoa(i), "data-"+strconv.Itoa(i))
	}

	wg.Wait()

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	for condition.Load() {
		time.Sleep(100 * time.Millisecond)
	}
}
