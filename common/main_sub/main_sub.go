package main_sub

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/log"
	"github.com/common-library/go/command-line/flags"
	"github.com/common-library/go/log/klog"
	"github.com/common-library/go/log/slog"
)

type ServerKind string

const (
	CloudEvents ServerKind = "cloudEvents"
	GRPC        ServerKind = "gRPC"
	Http        ServerKind = "http"
	LongPolling ServerKind = "longPolling"
	Socket      ServerKind = "socket"
)

type Main struct {
	serverKind ServerKind
}

func (this *Main) initialize(serverKind ServerKind) error {
	if err := this.parseFlag(); err != nil {
		return err
	} else if err := config.Read(flags.Get[string]("config-file")); err != nil {
		return err
	} else {
		if serverKind != CloudEvents {
			log.Initialize(string(serverKind))
		}

		return nil
	}
}

func (this *Main) parseFlag() error {
	flagInfos := []flags.FlagInfo{
		{FlagName: "config-file", Usage: "config/config.yaml", DefaultValue: string("")},
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

func (this *Main) RunWithKlog(serverKind ServerKind, start func() error, stop func() error) error {
	defer klog.Flush()

	if err := this.initialize(serverKind); err != nil {
		return err
	} else if err := start(); err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	return stop()
}

func (this *Main) RunWithSlog(serverKind ServerKind, start func(*slog.Log) error, stop func(*slog.Log) error) error {
	defer log.Log.Flush()

	if err := this.initialize(serverKind); err != nil {
		return err
	} else if err := start(&log.Log); err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	<-signals

	return stop(&log.Log)
}
