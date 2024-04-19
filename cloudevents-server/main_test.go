package main

import (
	"flag"
	"net/http"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/base-server/go/config"
	"github.com/common-library/go/event/cloudevents"
)

func TestMain1(t *testing.T) {
	os.Args = []string{"test"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if err := (&Main{}).Run(); err.Error() != "invalid flag" {
		t.Errorf("invalid error - (%#v)", err)
	}
}

func TestMain2(t *testing.T) {
	os.Args = []string{"test", "-config-file=invalid"}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	if err := (&Main{}).Run(); err.Error() != "open invalid: no such file or directory" {
		t.Errorf("invalid error - (%s)", err.Error())
	}
}

func TestMain3(t *testing.T) {
	path, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	configFile := path + "/../config/CloudEventsServer.config"

	wg := new(sync.WaitGroup)

	wg.Add(1)
	go func() {
		defer wg.Done()

		os.Args = []string{"test", "-config-file=" + configFile}
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

		main()
	}()
	time.Sleep(200 * time.Millisecond)

	if serverConfig, err := config.Get[config.CloudEventsServer](configFile); err != nil {
		t.Fatal(err)
	} else if client, err := cloudevents.NewHttp("http://"+serverConfig.Address, nil, nil); err != nil {
		t.Fatal(err)
	} else {
		const eventID = "id 01"
		const eventType = "type 01"
		const eventSource = "source/01"

		sendEvent := cloudevents.NewEvent()
		sendEvent.SetID(eventID)
		sendEvent.SetType(eventType)
		sendEvent.SetSource(eventSource)

		for i := 0; i < 100; i++ {
			if receiveEvent, result := client.Request(sendEvent); result.IsUndelivered() {
				t.Fatal(result.Error())
			} else if statusCode, err := result.GetHttpStatusCode(); err != nil {
				t.Fatal(err)
			} else if statusCode != http.StatusOK {
				t.Fatal("invalid -", statusCode)
			} else {
				if receiveEvent.ID() != eventID {
					t.Error("invalid -", receiveEvent.ID())
				} else if receiveEvent.Type() != eventType {
					t.Error("invalid -", receiveEvent.Type())
				} else if receiveEvent.Source() != eventSource {
					t.Error("invalid -", receiveEvent.Source())
				}
			}
		}
	}

	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Error(err)
	}

	wg.Wait()
}
