package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/base-server-go/socket-server/log"
	command_line_flag "github.com/heaven-chp/common-library-go/command-line/flag"
	"github.com/heaven-chp/common-library-go/socket"
)

type Main struct {
	server             socket.Server
	socketServerConfig config.SocketServer
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
		{FlagName: "config_file", Usage: "config/SocketServer.config", DefaultValue: string("")},
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

	if socketServerConfig, err := config.Get[config.SocketServer](fileName); err != nil {
		return err
	} else {
		this.socketServerConfig = socketServerConfig
		return nil
	}
}

func (this *Main) initializeLog() error {
	log.Initialize(this.socketServerConfig)

	return nil
}

func (this *Main) finalizeLog() error {
	log.Server.Flush()

	return nil
}

func (this *Main) initializeServer() error {
	acceptSuccessFunc := func(client socket.Client) {
		log.Server.Debug("start", "network", client.GetRemoteAddr().Network(), "address", client.GetRemoteAddr().String())
		log.Server.Debug("end", "network", client.GetRemoteAddr().Network(), "address", client.GetRemoteAddr().String())

		read := func(readJob func(readData string) error) error {
			readData, err := client.Read(1024)
			if err != nil {
				return err
			}

			log.Server.Debug("read", "data", readData)

			return readJob(readData)
		}

		write := func(writeData string) error {
			writeLen, err := client.Write(writeData)
			if err != nil {
				return err
			}
			if writeLen != len(writeData) {
				return errors.New(fmt.Sprintf("invalid write - (%d)(%d)", writeLen, len(writeData)))
			}

			log.Server.Debug("write", "data", writeData)

			return nil
		}

		err := write("greeting")
		if err != nil {
			log.Server.Error(err.Error())
			return
		}

		readJob := func(readData string) error {
			return write("[response] " + readData)
		}
		err = read(readJob)
		if err != nil {
			log.Server.Error(err.Error())
			return
		}
	}

	acceptFailureFunc := func(err error) {
		log.Server.Error(err.Error())
	}

	err := this.server.Start("tcp", this.socketServerConfig.Address, this.socketServerConfig.ClientPoolSize, acceptSuccessFunc, acceptFailureFunc)
	if err != nil {
		panic(err)
	}

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
