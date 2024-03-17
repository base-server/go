package main

import (
	"errors"
	"flag"
	net_http "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/base-server-go/http-server/handler"
	"github.com/heaven-chp/base-server-go/http-server/log"
	"github.com/heaven-chp/base-server-go/http-server/swagger_docs"
	command_line_flag "github.com/heaven-chp/common-library-go/command-line/flag"
	"github.com/heaven-chp/common-library-go/http"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Main struct {
	server           http.Server
	httpServerConfig config.HttpServer
}

func (this *Main) initialize() error {
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

	err = this.initializeSwagger()
	if err != nil {
		return err
	}

	err = this.initializeServer()
	if err != nil {
		return err
	}

	return nil
}

func (this *Main) finalize() error {
	defer this.finalizeLog()

	return this.finalizeServer()
}

func (this *Main) initializeFlag() error {
	err := command_line_flag.Parse([]command_line_flag.FlagInfo{
		{FlagName: "config_file", Usage: "config/HttpServer.config", DefaultValue: string("")},
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

	if httpServerConfig, err := config.Get[config.HttpServer](fileName); err != nil {
		return err
	} else {
		this.httpServerConfig = httpServerConfig
		return nil
	}
}

func (this *Main) initializeLog() error {
	log.Initialize(this.httpServerConfig)

	return nil
}

func (this *Main) initializeSwagger() error {
	swagger_docs.SwaggerInfo.Version = "1.0"
	swagger_docs.SwaggerInfo.Host = this.httpServerConfig.SwaggerAddress
	swagger_docs.SwaggerInfo.BasePath = ""
	swagger_docs.SwaggerInfo.Title = "http server"
	swagger_docs.SwaggerInfo.Description = ""

	return nil
}

func (this *Main) initializeServer() error {
	this.server.AddPathPrefixHandler(this.httpServerConfig.SwaggerUri, httpSwagger.WrapHandler)

	this.server.AddHandler("/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", net_http.MethodGet, handler.Get)
	this.server.AddHandler("/v1/test", net_http.MethodPost, handler.Post)
	this.server.AddHandler("/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", net_http.MethodDelete, handler.Delete)

	return this.server.Start(this.httpServerConfig.ServerAddress, func(err error) { log.Server.Error(err.Error()) })
}

func (this *Main) finalizeLog() error {
	log.Server.Flush()
	return nil
}

func (this *Main) finalizeServer() error {
	return this.server.Stop(this.httpServerConfig.ShutdownTimeout)
}

func (this *Main) Run() error {
	err := this.initialize()
	if err != nil {
		return err
	}
	defer this.finalize()

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
