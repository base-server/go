package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/base-server/go/config"
	"github.com/base-server/go/grpc-server/log"
	"github.com/common-library/go/command-line/flags"
	"github.com/common-library/go/event/cloudevents"
	"github.com/common-library/go/log/klog"
)

type Main struct {
	server                  cloudevents.Server
	cloudEventsServerConfig config.CloudEventsServer
}

func (this *Main) initialize() error {
	if err := this.parseFlag(); err != nil {
		return err
	} else if err := this.setConfig(); err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Main) parseFlag() error {
	flagInfos := []flags.FlagInfo{
		{FlagName: "config-file", Usage: "config/CloudEventsServer.config", DefaultValue: string("")},
	}

	if err := flags.Parse(flagInfos); err != nil {
		flag.Usage()
		return err
	} else if flag.NFlag() != 1 {
		flag.Usage()
		return errors.New("invalid flag")
	} else {
		return nil
	}
}

func (this *Main) setConfig() error {
	fileName := flags.Get[string]("config-file")

	if cloudEventsServerConfig, err := config.Get[config.CloudEventsServer](fileName); err != nil {
		return err
	} else {
		this.cloudEventsServerConfig = cloudEventsServerConfig
		return nil
	}
}

func (this *Main) Run() error {
	if err := this.initialize(); err != nil {
		return err
	}

	address := this.cloudEventsServerConfig.Address
	handler := func(event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
		klog.InfoS("handler", "event", event.String())

		responseEvent := event.Clone()
		return &responseEvent, cloudevents.NewHTTPResult(http.StatusOK, "")
	}

	failureFunc := func(err error) { klog.ErrorS(err, "") }

	if err := this.server.Start(address, handler, failureFunc); err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	log.Server.Info("signal", "kind", <-signals)

	shutdownTimeout := this.cloudEventsServerConfig.ShutdownTimeout
	if duration, err := time.ParseDuration(shutdownTimeout); err != nil {
		return err
	} else {
		return this.server.Stop(duration)
	}
}

func main() {
	defer klog.Flush()

	klog.InfoS("main start")
	defer klog.InfoS("main end")

	if err := (&Main{}).Run(); err != nil {
		klog.ErrorS(err, "")
	}
}
