package main

import (
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/common-library-go/log"
	"github.com/heaven-chp/common-library-go/socket"
)

type Main struct {
	configFile string

	socketServerConfig config.SocketServer

	server socket.Server
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
	configFile := flag.String("config_file", "", "config file")
	flag.Parse()

	if flag.NFlag() != 1 {
		flag.Usage()
		return errors.New("invalid flag")
	}

	this.configFile = *configFile

	return nil
}

func (this *Main) initializeConfig() error {
	return config.Parsing(&this.socketServerConfig, this.configFile)
}

func (this *Main) initializeLog() error {
	level, err := log.ToIntLevel(this.socketServerConfig.LogLevel)
	if err != nil {
		return err
	}

	return log.Initialize(level, this.socketServerConfig.LogOutputPath, this.socketServerConfig.LogFileNamePrefix)
}

func (this *Main) finalizeLog() error {
	return log.Finalize()
}

func (this *Main) initializeServer() error {
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
		err := this.server.Start("tcp", this.socketServerConfig.Address, this.socketServerConfig.ClientPoolSize, serverJob)
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
