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

	return main.finalizeServer()
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
	serverJob := func(client socket.Client) {
		read := func(readJob func(readData string) bool) bool {
			readData, err := client.Read(1024)
			if err != nil {
				log.Error("read fail - error : (%s)", err.Error)
				return false
			}

			return readJob(readData)
		}

		write := func(writeData string) bool {
			writeLen, err := client.Write(writeData)
			if err != nil {
				log.Error("write error - (%s)", err.Error())
				return false
			}
			if writeLen != len(writeData) {
				log.Error("invalid write - (%d)(%d)", writeLen, len(writeData))
				return false
			}

			log.Debug("write data - data : (%s)", writeData)

			return true
		}

		if write("greeting") == false {
			return
		}

		readJob1 := func(readData string) bool {
			log.Debug("read data - data : (%s)", readData)

			if write("[response] "+readData) == false {
				return false
			}

			return true
		}
		if read(readJob1) == false {
			return
		}
	}

	go func() {
		err := main.server.Start("tcp", main.socketServerConfig.Address, main.socketServerConfig.ClientPoolSize, serverJob)
		if err != nil {
			panic(err)
		}
	}()

	return nil
}

func (main *Main) finalizeServer() error {
	return main.server.Stop()
}

func (main *Main) Run() error {
	err := main.Initialize()
	if err != nil {
		return err
	}
	defer main.Finalize()

	log.Info("process start")
	defer log.Info("process end")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info("signal : (%s)", <-signals)

	return nil
}

func main() {
	main := Main{}

	err := main.Run()
	if err != nil {
		log.Error(err.Error())
	}
}
