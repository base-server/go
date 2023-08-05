package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	command_line_argument "github.com/heaven-chp/common-library-go/command-line-argument"
	log "github.com/heaven-chp/common-library-go/log/file"
	long_polling "github.com/heaven-chp/common-library-go/long-polling"
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
	server                  long_polling.Server
	longPollingServerConfig config.LongPollingServer
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
		{FlagName: "config_file", Usage: "config/LongPollingServer.config", DefaultValue: string("")},
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
	return config.Parsing(&this.longPollingServerConfig, command_line_argument.Get("config_file").(string))
}

func (this *Main) initializeLog() error {
	return log_instance().Initialize(log.Setting{
		Level:           this.longPollingServerConfig.Log.Level,
		OutputPath:      this.longPollingServerConfig.Log.OutputPath,
		FileNamePrefix:  this.longPollingServerConfig.Log.FileNamePrefix,
		PrintCallerInfo: this.longPollingServerConfig.Log.PrintCallerInfo,
		ChannelSize:     this.longPollingServerConfig.Log.ChannelSize})
}

func (this *Main) finalizeLog() error {
	return log_instance().Finalize()
}

func (this *Main) initializeServer() error {
	serverInfo := long_polling.ServerInfo{
		Address:                        this.longPollingServerConfig.Address,
		Timeout:                        this.longPollingServerConfig.Timeout,
		SubscriptionURI:                this.longPollingServerConfig.SubscriptionURI,
		HandlerToRunBeforeSubscription: func(w http.ResponseWriter, r *http.Request) bool { return true },
		PublishURI:                     this.longPollingServerConfig.PublishURI,
		HandlerToRunBeforePublish:      func(w http.ResponseWriter, r *http.Request) bool { return true }}

	filePersistorInfo := long_polling.FilePersistorInfo{
		Use:                     this.longPollingServerConfig.FilePersistorInfo.Use,
		FileName:                this.longPollingServerConfig.FilePersistorInfo.FileName,
		WriteBufferSize:         this.longPollingServerConfig.FilePersistorInfo.WriteBufferSize,
		WriteFlushPeriodSeconds: this.longPollingServerConfig.FilePersistorInfo.WriteFlushPeriodSeconds}

	err := this.server.Start(serverInfo, filePersistorInfo, func(err error) { log_instance().Error(err) })
	if err != nil {
		panic(err)
	}

	return nil
}

func (this *Main) finalizeServer() error {
	return this.server.Stop(this.longPollingServerConfig.ShutdownTimeout)
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
