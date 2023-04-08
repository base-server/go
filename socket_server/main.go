package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/common-library-go/json"
	"github.com/heaven-chp/common-library-go/log"
	"github.com/heaven-chp/common-library-go/socket"
)

type Main struct {
	configFile string

	socketServerConfig config.SocketServer

	server socket.Server
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

func (main *Main) Finalize() error {
	defer main.finalizeLog()

	main.finalizeServer()

	return nil
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
	return json.ToStructFromFile(main.configFile, &main.socketServerConfig)
}

func (main *Main) initializeLog() error {
	level, err := log.ToIntLevel(main.socketServerConfig.LogLevel)
	if err != nil {
		return err
	}

	return log.Initialize(level, main.socketServerConfig.LogOutputPath, main.socketServerConfig.LogFileNamePrefix)
}

func (main *Main) finalizeLog() error {
	return log.Finalize()
}

func (main *Main) initializeServer() error {
	jobFunc := func(client socket.Client) {
		client.Write("greeting")

		readData, err := client.Read(1024)
		if err != nil {
			log.Error("read fail - error : (%s)", err.Error)
			return
		}

		log.Debug("read data - data : (%s)", readData)

		writeData := "[response] " + readData
		writeLen, err := client.Write(writeData)
		if err != nil {
			log.Error("write fail - error : (%s)", err.Error)
			return
		}
		if writeLen != len(writeData) {
			log.Error("write len is different - writeLen : (%d), len(writeData) : (%d)", writeLen, len(writeData))
		}

		log.Debug("write data - data : (%s)", writeData)
	}

	err := main.server.Initialize("tcp", main.socketServerConfig.Address, main.socketServerConfig.ClientPoolSize, jobFunc)
	if err != nil {
		return err
	}

	return nil
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
	main.Run()
	log.Info("process end")
}
