package main

import (
	"errors"
	"flag"
	net_http "net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/base-server-go/http-server/swagger_docs"
	command_line_argument "github.com/heaven-chp/common-library-go/command-line-argument"
	"github.com/heaven-chp/common-library-go/http"
	"github.com/heaven-chp/common-library-go/log"
	httpSwagger "github.com/swaggo/http-swagger"
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
	err := command_line_argument.Set([]command_line_argument.CommandLineArgumentInfo{
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
	return config.Parsing(&this.httpServerConfig, command_line_argument.Get("config_file").(string))
}

func (this *Main) initializeLog() error {
	level, err := log.ToIntLevel(this.httpServerConfig.LogLevel)
	if err != nil {
		return err
	}

	return log.Initialize(level, this.httpServerConfig.LogOutputPath, this.httpServerConfig.LogFileNamePrefix)
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

	this.server.AddHandler("/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", net_http.MethodGet, testGet)
	this.server.AddHandler("/v1/test", net_http.MethodPost, testPost)
	this.server.AddHandler("/v1/test/{id:[a-z,A-Z][a-z,A-Z,0-9,--,_,.]+}", net_http.MethodDelete, testDelete)

	return this.server.Start(this.httpServerConfig.ServerAddress, func(err error) { log.Error(err.Error()) })
}

func (this *Main) finalizeLog() error {
	return log.Finalize()
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

	log.Info("process start")
	defer log.Info("process end")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	log.Info("signal : (%s)", <-signals)

	return nil
}

func main() {
	main := Main{}
	err := main.Run()
	if err != nil {
		log.Error(err.Error())
	}
}
