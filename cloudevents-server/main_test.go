package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/common-library/go/event/cloudevents"
)

func TestMain(t *testing.T) {
	const configFile = "../common/config/config.yaml"

	wg := new(sync.WaitGroup)
	wg.Go(func() {

		os.Args = []string{"test", "-config-file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()
	})
	time.Sleep(200 * time.Millisecond)

	if err := config.Read(configFile); err != nil {
		t.Fatal(err)
	} else if client, err := cloudevents.NewHttp("http://"+config.Get("cloudEvents.address").(string), nil, nil); err != nil {
		t.Fatal(err)
	} else {
		const eventID = "id 01"
		const eventType = "type 01"
		const eventSource = "source/01"

		sendEvent := cloudevents.NewEvent()
		sendEvent.SetID(eventID)
		sendEvent.SetType(eventType)
		sendEvent.SetSource(eventSource)

		for range 100 {
			if receiveEvent, result := client.Request(sendEvent); result.IsUndelivered() {
				t.Fatal(result.Error())
			} else if statusCode, err := result.GetHttpStatusCode(); err != nil {
				t.Fatal(err)
			} else if statusCode != http.StatusOK {
				t.Fatal(statusCode)
			} else {
				if receiveEvent.ID() != eventID {
					t.Fatal(receiveEvent.ID(), ",", eventID)
				} else if receiveEvent.Type() != eventType {
					t.Fatal(receiveEvent.Type(), ",", eventType)
				} else if receiveEvent.Source() != eventSource {
					t.Fatal(receiveEvent.Source(), ",", eventSource)
				}
			}
		}
	}

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}

	wg.Wait()
}
