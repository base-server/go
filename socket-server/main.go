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
	"github.com/heaven-chp/common-library-go/command-line/flags"
	"github.com/heaven-chp/common-library-go/socket"
)

type Main struct {
	server             socket.Server
	socketServerConfig config.SocketServer
}

func (this *Main) initialize() error {
	if err := this.parseFlag(); err != nil {
		return err
	} else if err := this.setConfig(); err != nil {
		return err
	} else {
		log.Initialize(this.socketServerConfig)
		return nil
	}
}

func (this *Main) parseFlag() error {
	flagInfos := []flags.FlagInfo{
		{FlagName: "config_file", Usage: "config/SocketServer.config", DefaultValue: string("")},
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

	if socketServerConfig, err := config.Get[config.SocketServer](fileName); err != nil {
		return err
	} else {
		this.socketServerConfig = socketServerConfig
		return nil
	}
}

func (this *Main) startServer() error {
	acceptSuccessFunc := func(client socket.Client) {
		log.Server.Debug("start", "network", client.GetRemoteAddr().Network(), "address", client.GetRemoteAddr().String())
		log.Server.Debug("end", "network", client.GetRemoteAddr().Network(), "address", client.GetRemoteAddr().String())

		read := func(readJob func(readData string) error) error {
			if readData, err := client.Read(1024); err != nil {
				return err
			} else {
				log.Server.Debug("read", "data", readData)

				return readJob(readData)
			}
		}

		write := func(writeData string) error {
			if writeLen, err := client.Write(writeData); err != nil {
				return err
			} else if writeLen != len(writeData) {
				return errors.New(fmt.Sprintf("invalid write - (%d)(%d)", writeLen, len(writeData)))
			} else {
				log.Server.Debug("write", "data", writeData)

				return nil
			}
		}

		if err := write("greeting"); err != nil {
			log.Server.Error(err.Error())
			return
		}

		readJob := func(readData string) error {
			return write("[response] " + readData)
		}
		if err := read(readJob); err != nil {
			log.Server.Error(err.Error())
			return
		}
	}

	acceptFailureFunc := func(err error) {
		log.Server.Error(err.Error())
	}

	return this.server.Start("tcp", this.socketServerConfig.Address, this.socketServerConfig.ClientPoolSize, acceptSuccessFunc, acceptFailureFunc)
}

func (this *Main) Run() error {
	defer log.Server.Flush()

	if err := this.initialize(); err != nil {
		return err
	}

	log.Server.Info("process start")
	defer log.Server.Info("process end")

	if err := this.startServer(); err != nil {
		return err
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Server.Info("signal", "kind", <-signals)

	return this.server.Stop()
}

func main() {
	if err := (&Main{}).Run(); err != nil {
		log.Server.Error(err.Error())
		log.Server.Flush()
	}
}
