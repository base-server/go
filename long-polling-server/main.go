package main

import (
	"net/http"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/common-library/go/log/slog"
	long_polling "github.com/common-library/go/long-polling"
)

func main() {
	server := long_polling.Server{}

	start := func(log *slog.Log) error {
		serverInfo := long_polling.ServerInfo{
			Address:                        config.Get("longPolling.address").(string),
			TimeoutSeconds:                 config.Get("longPolling.timeoutSeconds").(int),
			SubscriptionURI:                config.Get("longPolling.subscriptionURI").(string),
			HandlerToRunBeforeSubscription: func(w http.ResponseWriter, r *http.Request) bool { return true },
			PublishURI:                     config.Get("longPolling.publishURI").(string),
			HandlerToRunBeforePublish:      func(w http.ResponseWriter, r *http.Request) bool { return true }}

		filePersistorInfo := long_polling.FilePersistorInfo{
			Use:                     config.Get("longPolling.filePersistorInfo.use").(bool),
			FileName:                config.Get("longPolling.filePersistorInfo.fileName").(string),
			WriteBufferSize:         config.Get("longPolling.filePersistorInfo.writeBufferSize").(int),
			WriteFlushPeriodSeconds: config.Get("longPolling.filePersistorInfo.writeFlushPeriodSeconds").(int)}

		return server.Start(serverInfo, filePersistorInfo, func(err error) { log.Error(err.Error()) })
	}

	stop := func(log *slog.Log) error {
		shutdownTimeout := config.Get("longPolling.shutdownTimeout").(string)
		if duration, err := time.ParseDuration(shutdownTimeout); err != nil {
			return err
		} else {
			return server.Stop(duration)
		}
	}

	if err := (&main_sub.Main{}).RunWithSlog(main_sub.LongPolling, start, stop); err != nil {
		panic(err)
	}
}
