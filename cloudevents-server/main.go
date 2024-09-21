package main

import (
	"net/http"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/common-library/go/event/cloudevents"
	"github.com/common-library/go/log/klog"
)

func main() {

	server := cloudevents.Server{}

	start := func() error {
		address := config.Get("cloudEvents.address").(string)
		handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
			klog.InfoS("handler", "event", event.String())

			responseEvent := event.Clone()
			return &responseEvent, cloudevents.NewHTTPResult(http.StatusOK, "")
		}
		failureFunc := func(err error) { klog.ErrorS(err, "") }

		return server.Start(address, handler, failureFunc)
	}

	stop := func() error {
		shutdownTimeout := config.Get("cloudEvents.shutdownTimeout").(string)
		if duration, err := time.ParseDuration(shutdownTimeout); err != nil {
			return err
		} else {
			return server.Stop(duration)
		}
	}

	if err := (&main_sub.Main{}).RunWithKlog(main_sub.CloudEvents, start, stop); err != nil {
		panic(err)
	}
}
