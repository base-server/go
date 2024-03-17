package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/base-server-go/grpc-server/log"
	command_line_flag "github.com/heaven-chp/common-library-go/command-line/flag"
	"github.com/heaven-chp/common-library-go/grpc"
	"github.com/heaven-chp/common-library-go/grpc/sample"
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
	flagInfos := []command_line_flag.FlagInfo{
		{FlagName: "config_file", Usage: "config/GrpcServer.config", DefaultValue: string("")},
	}

	if err := command_line_flag.Parse(flagInfos); err != nil {
		return nil
	} else if flag.NFlag() != 1 {
		flag.Usage()
		return errors.New("invalid flag")
	} else {
		return nil
	}
}

func (this *Main) setConfig() error {
	fileName := command_line_flag.Get[string]("config_file")

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
