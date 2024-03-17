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

func (this *Main) Initialize() error {
	err := this.initializeFlag()
	if err != nil {
		return err
	}

	err = this.initializeConfig()
	if err != nil {
		return err
	}

	err = this.initializeLog()
	if err != nil {
		return err
	}

	err = this.initializeServer()
	if err != nil {
		return err
	}

	return nil
}

func (this *Main) Finalize() error {
	defer this.finalizeLog()

	return this.finalizeServer()
}

func (this *Main) initializeFlag() error {
	err := command_line_flag.Parse([]command_line_flag.FlagInfo{
		{FlagName: "config_file", Usage: "config/GrpcServer.config", DefaultValue: string("")},
	})
	if err != nil {
		return nil
	}

	if flag.NFlag() != 1 {
		flag.Usage()
		return errors.New("invalid flag")
	}

	return nil
}

func (this *Main) initializeConfig() error {
	fileName := command_line_flag.Get[string]("config_file")

	if grpcServerConfig, err := config.Get[config.GrpcServer](fileName); err != nil {
		return err
	} else {
		this.grpcServerConfig = grpcServerConfig
		return nil
	}
}

func (this *Main) initializeLog() error {
	log.Initialize(this.grpcServerConfig)

	return nil
}

func (this *Main) finalizeLog() error {
	log.Server.Flush()

	return nil
}

func (this *Main) initializeServer() error {
	go func() {
		err := this.server.Start(this.grpcServerConfig.Address, &sample.Server{})
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (this *Main) finalizeServer() error {
	return this.server.Stop()
}

func (this *Main) Run() error {
	err := this.Initialize()
	if err != nil {
		return err
	}
	defer this.Finalize()

	log.Server.Info("process start")
	defer log.Server.Info("process end")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Server.Info("signal", "kind", <-signals)

	return nil
}

func main() {
	main := Main{}

	err := main.Run()
	if err != nil {
		log.Server.Error(err.Error())
	}
}
