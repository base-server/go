package main

import (
	net_http "net/http"
	"time"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/base-server/go/http-server/handler"
	"github.com/base-server/go/http-server/swagger_docs"
	"github.com/common-library/go/http"
	"github.com/common-library/go/log/slog"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func setSwaggerInfo() {
	swagger_docs.SwaggerInfo.Version = "1.0"
	swagger_docs.SwaggerInfo.Host = config.Get("http.swaggerAddress").(string)
	swagger_docs.SwaggerInfo.BasePath = ""
	swagger_docs.SwaggerInfo.Title = "http server"
	swagger_docs.SwaggerInfo.Description = ""
}

func setHandler(server *http.Server) {
	server.RegisterPathPrefixHandler(config.Get("http.swaggerUri").(string), httpSwagger.WrapHandler)

	server.RegisterHandlerFunc("/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", net_http.MethodGet, handler.Get)
	server.RegisterHandlerFunc("/v1/test", net_http.MethodPost, handler.Post)
	server.RegisterHandlerFunc("/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", net_http.MethodDelete, handler.Delete)
}

func main() {
	setSwaggerInfo()

	server := http.Server{}
	setHandler(&server)

	start := func(log *slog.Log) error {
		listenAndServeFailureFunc := func(err error) { log.Error(err.Error()) }
		return server.Start(config.Get("http.serverAddress").(string), listenAndServeFailureFunc)
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
