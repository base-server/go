package main

import (
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/base-server/go/config"
	"github.com/base-server/go/long-polling-server/log"
	"github.com/common-library/go/command-line/flags"
	long_polling "github.com/common-library/go/long-polling"
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
	flagInfos := []flags.FlagInfo{
		{FlagName: "config_file", Usage: "config/LongPollingServer.config", DefaultValue: string("")},
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
		Timeout:                        this.longPollingServerConfig.TimeoutSeconds,
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
	shutdownTimeout := this.longPollingServerConfig.ShutdownTimeout
	if duration, err := time.ParseDuration(shutdownTimeout); err != nil {
		return err
	} else {
		return this.server.Stop(duration)
	}
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
