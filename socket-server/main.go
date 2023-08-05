package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	command_line_argument "github.com/heaven-chp/common-library-go/command-line-argument"
	log "github.com/heaven-chp/common-library-go/log/file"
	"github.com/heaven-chp/common-library-go/socket"
	"github.com/heaven-chp/common-library-go/utility"
)

var onceForLog sync.Once
var fileLog *log.FileLog

func log_instance() *log.FileLog {
	onceForLog.Do(func() {
		fileLog = &log.FileLog{}
	})

	return fileLog
}

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
	err := command_line_argument.Set([]command_line_argument.CommandLineArgumentInfo{
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
	return config.Parsing(&this.socketServerConfig, command_line_argument.Get("config_file").(string))
}

func (this *Main) initializeLog() error {
	return log_instance().Initialize(log.Setting{
		Level:           this.socketServerConfig.Log.Level,
		OutputPath:      this.socketServerConfig.Log.OutputPath,
		FileNamePrefix:  this.socketServerConfig.Log.FileNamePrefix,
		PrintCallerInfo: this.socketServerConfig.Log.PrintCallerInfo,
		ChannelSize:     this.socketServerConfig.Log.ChannelSize})
}

func (this *Main) finalizeLog() error {
	return log_instance().Finalize()
}

func (this *Main) initializeServer() error {
	acceptSuccessFunc := func(client socket.Client) {
		prefixLog := ""
		callerInfo, err := utility.GetCallerInfo(1)
		if err != nil {
			log_instance().Error(err)
		} else {
			prefixLog = "[goroutine-id:" + strconv.Itoa(callerInfo.GoroutineID) + "] : "
		}

		log_instance().Debugf("%sstart - (%s)(%s)", prefixLog, client.GetRemoteAddr().Network(), client.GetRemoteAddr().String())
		defer log_instance().Debugf("%send - (%s)(%s)", prefixLog, client.GetRemoteAddr().Network(), client.GetRemoteAddr().String())

		read := func(readJob func(readData string) error) error {
			readData, err := client.Read(1024)
			if err != nil {
				return err
			}

			log_instance().Debugf("%sread data - data : (%s)", prefixLog, readData)

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

			log_instance().Debugf("%swrite data - data : (%s)", prefixLog, writeData)

			return nil
		}

		err = write("greeting")
		if err != nil {
			log_instance().Errorf("%s%s", prefixLog, err.Error())
			return
		}

		readJob := func(readData string) error {
			return write("[response] " + readData)
		}
		err = read(readJob)
		if err != nil {
			log_instance().Errorf("%s%s", prefixLog, err.Error())
			return
		}
	}

	acceptFailureFunc := func(err error) {
		log_instance().Error(err)
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

	log_instance().Info("process start")
	defer log_instance().Info("process end")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log_instance().Infof("signal : (%s)", <-signals)

	return nil
}

func main() {
	main := Main{}

	err := main.Run()
	if err != nil {
		log_instance().Error(err)
	}
}
