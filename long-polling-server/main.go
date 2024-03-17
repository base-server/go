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

func (this *Main) initialize() error {
	if err := this.parseFlag(); err != nil {
		return err
	} else if err := this.setConfig(); err != nil {
		return err
	} else {
		log.Initialize(this.longPollingServerConfig)

		return nil
	}
}

func (this *Main) parseFlag() error {
	flagInfos := []command_line_flag.FlagInfo{
		{FlagName: "config_file", Usage: "config/LongPollingServer.config", DefaultValue: string("")},
	}

	if err := command_line_flag.Parse(flagInfos); err != nil {
		return nil
	} else if flag.NFlag() != 1 {
		flag.Usage()
		return errors.New("invalid flag")
	} else {
		return nil
	}
}

func (this *Main) setConfig() error {
	fileName := command_line_flag.Get[string]("config_file")

	if longPollingServerConfig, err := config.Get[config.LongPollingServer](fileName); err != nil {
		return err
	} else {
		this.longPollingServerConfig = longPollingServerConfig
		return nil
	}
}

func (this *Main) startServer() error {
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

	return this.server.Start(serverInfo, filePersistorInfo, func(err error) { log.Server.Error(err.Error()) })
}

func (this *Main) stopServer() error {
	return this.server.Stop(this.longPollingServerConfig.ShutdownTimeout)
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

	return this.stopServer()
}

func main() {
	if err := (&Main{}).Run(); err != nil {
		log.Server.Error(err.Error())
		log.Server.Flush()
	}
}
