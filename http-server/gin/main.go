package main

import (
	"net/http"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/base-server/go/http-server/handler"
	"github.com/common-library/go/http/gin"
	"github.com/common-library/go/log/slog"
)

func setHandler(server *gin.Server) {
	server.RegisterHandler(http.MethodGet, "/v1/test/:id", gin.WrapHandlerFunc(handler.Get))
	server.RegisterHandler(http.MethodPost, "/v1/test", gin.WrapHandlerFunc(handler.Post))
	server.RegisterHandler(http.MethodDelete, "/v1/test/:id", gin.WrapHandlerFunc(handler.Delete))
}

func main() {
	server := gin.Server{}
	setHandler(&server)

	start := func(log *slog.Log) error {
		listenAndServeFailureFunc := func(err error) { log.Error(err.Error()) }
		return server.Start(config.Get("http.gin.address").(string), listenAndServeFailureFunc)
	}

	stop := func(log *slog.Log) error {
		shutdownTimeout := config.Get("http.shutdownTimeout").(string)
		if duration, err := time.ParseDuration(shutdownTimeout); err != nil {
			return err
		} else {
			return server.Stop(duration)
		}
	}

	if err := (&main_sub.Main{}).RunWithSlog(main_sub.Http, start, stop); err != nil {
		panic(err)
	}
}
