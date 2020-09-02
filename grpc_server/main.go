package main

import (
	"errors"
	"flag"
	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/common-library-go/grpc"
	"github.com/heaven-chp/common-library-go/grpc/sample"
	"github.com/heaven-chp/common-library-go/json"
	"github.com/heaven-chp/common-library-go/log"
	"os"
	"os/signal"
	"syscall"
)

type Main struct {
	configFile string

	grpcServerConfig config.GrpcServer

	server grpc.Server
}

func (main *Main) Initialize() error {
	err := main.initializeFlag()
	if err != nil {
		return err
	}

	err = main.initializeConfig()
	if err != nil {
		return err
	}

	err = main.initializeLog()
	if err != nil {
		return err
	}

	err = main.initializeServer()
	if err != nil {
		return err
	}

	return nil
}

func (main *Main) Finalize() {
	defer main.finalizeLog()

	main.finalizeServer()
}

func (main *Main) initializeFlag() error {
	configFile := flag.String("config_file", "", "config file")
	flag.Parse()

	if flag.NFlag() != 1 {
		flag.Usage()
		return errors.New("invalid flag")
	}

	main.configFile = *configFile

	return nil
}

func (main *Main) initializeConfig() error {
	return json.ToStructFromFile(main.configFile, &main.grpcServerConfig)
}

func (main *Main) initializeLog() error {
	level, err := log.ToIntLevel(main.grpcServerConfig.LogLevel)
	if err != nil {
		return err
	}

	return log.Initialize(level, main.grpcServerConfig.LogOutputPath, main.grpcServerConfig.LogFileNamePrefix)
}

func (main *Main) finalizeLog() error {
	return log.Finalize()
}

func (main *Main) initializeServer() error {
	return main.server.Initialize(main.grpcServerConfig.Address, &sample.Server{})
}

func (main *Main) finalizeServer() error {
	return main.server.Finalize()
}

func (main *Main) Run() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go main.server.Run()

	log.Info("signal : (%s)", <-signals)
}

func main() {
	var main Main
	err := main.Initialize()
	if err != nil {
		log.Error("Initialize fail - error : (%s)", err.Error())
		return
	}
	defer main.Finalize()

	log.Info("process start")
	defer log.Info("process end")

	main.Run()
}
