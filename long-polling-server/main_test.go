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

	"github.com/base-server/go/common/config"
	"github.com/common-library/go/file"
	long_polling "github.com/common-library/go/long-polling"
)

func subscription(t *testing.T, request long_polling.SubscriptionRequest, count int, data string) (int64, string) {
	response, err := long_polling.Subscription("http://"+config.Get("longPolling.address").(string)+config.Get("longPolling.subscriptionURI").(string), nil, request, "", "", nil)
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

func publish(t *testing.T, category, data string) {
	request := long_polling.PublishRequest{Category: category, Data: data}
	response, err := long_polling.Publish("http://"+config.Get("longPolling.address").(string)+config.Get("longPolling.publishURI").(string), 10, nil, request, "", "", nil)
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

func TestMain(t *testing.T) {
	const configFile = "../common/config/config.yaml"

	// Load config file first
	if err := config.Read(configFile); err != nil {
		t.Fatal(err)
	}

	condition := atomic.Bool{}
	condition.Store(false)
	go func() {
		os.Args = []string{"test", "-config-file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		condition.Store(true)
		main()
		condition.Store(false)
	}()
	time.Sleep(500 * time.Millisecond)

	wg := new(sync.WaitGroup)

	clientJob := func(category, data string) {
		println(category, data)
		defer wg.Done()

		publish(t, category, data)
		timestamp, id := subscription(t, long_polling.SubscriptionRequest{Category: category, TimeoutSeconds: 300, SinceTime: 1}, 1, data)

		publish(t, category, data)
		publish(t, category, data)
		subscription(t, long_polling.SubscriptionRequest{Category: category, TimeoutSeconds: 300, SinceTime: timestamp, LastID: id}, 2, data)
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

	file.Remove(config.Get("longPolling.log.file.name").(string) + "." + config.Get("longPolling.log.file.extensionName").(string))
}
