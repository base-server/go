package main

import (
	"net/http"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/base-server/go/http-server/handler"
	"github.com/common-library/go/http/mux"
	"github.com/common-library/go/log/slog"
)

func setHandler(server *mux.Server) {
	server.RegisterHandlerFunc(http.MethodGet, "/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", http.HandlerFunc(handler.Get))
	server.RegisterHandlerFunc(http.MethodPost, "/v1/test", http.HandlerFunc(handler.Post))
	server.RegisterHandlerFunc(http.MethodDelete, "/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", http.HandlerFunc(handler.Delete))
}

func main() {
	server := mux.Server{}
	setHandler(&server)

	start := func(log *slog.Log) error {
		listenAndServeFailureFunc := func(err error) { log.Error(err.Error()) }
		return server.Start(config.Get("http.mux.address").(string), listenAndServeFailureFunc)
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
