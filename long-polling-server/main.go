package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/base-server-go/long-polling-server/log"
	command_line_flag "github.com/heaven-chp/common-library-go/command-line/flag"
	long_polling "github.com/heaven-chp/common-library-go/long-polling"
)

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
	err := command_line_flag.Parse([]command_line_flag.FlagInfo{
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
	fileName := command_line_flag.Get[string]("config_file")

	if longPollingServerConfig, err := config.Get[config.LongPollingServer](fileName); err != nil {
		return err
	} else {
		this.longPollingServerConfig = longPollingServerConfig
		return nil
	}
}

func (this *Main) initializeLog() error {
	log.Initialize(this.longPollingServerConfig)

	return nil
}

func (this *Main) finalizeLog() error {
	log.Server.Flush()

	return nil
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

	err := this.server.Start(serverInfo, filePersistorInfo, func(err error) { log.Server.Error(err.Error()) })
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
