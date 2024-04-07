package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/base-server/go/config"
	"github.com/base-server/go/grpc-server/log"
	"github.com/common-library/go/command-line/flags"
	"github.com/common-library/go/grpc"
	"github.com/common-library/go/grpc/sample"
)

type Main struct {
	server           grpc.Server
	grpcServerConfig config.GrpcServer
}

func (this *Main) initialize() error {
	if err := this.parseFlag(); err != nil {
		return err
	} else if err := this.setConfig(); err != nil {
		return err
	} else {
		log.Initialize(this.grpcServerConfig)

		return nil
	}
}

func (this *Main) parseFlag() error {
	flagInfos := []flags.FlagInfo{
		{FlagName: "config_file", Usage: "config/GrpcServer.config", DefaultValue: string("")},
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
	fileName := flags.Get[string]("config_file")

	if grpcServerConfig, err := config.Get[config.GrpcServer](fileName); err != nil {
		return err
	} else {
		this.grpcServerConfig = grpcServerConfig
		return nil
	}
}

func (this *Main) Run() error {
	defer log.Server.Flush()

	if err := this.initialize(); err != nil {
		return err
	}

	log.Server.Info("process start")
	defer log.Server.Info("process end")

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		log.Server.Info("signal", "kind", <-signals)

		this.server.Stop()
	}()

	return this.server.Start(this.grpcServerConfig.Address, &sample.Server{})
}

func main() {
	if err := (&Main{}).Run(); err != nil {
		log.Server.Error(err.Error())
		log.Server.Flush()
	}
}
