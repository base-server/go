package log

import (
	"strings"

	"github.com/heaven-chp/base-server-go/config"
	"github.com/heaven-chp/common-library-go/log/slog"
)

var Server slog.Log

func Initialize(serverConfig config.GrpcServer) {
	switch strings.ToLower(serverConfig.Log.Level) {
	case "trace":
		Server.SetLevel(slog.LevelTrace)
	case "debug":
		Server.SetLevel(slog.LevelDebug)
	case "info":
		Server.SetLevel(slog.LevelInfo)
	case "warn":
		Server.SetLevel(slog.LevelWarn)
	case "error":
		Server.SetLevel(slog.LevelError)
	case "fatal":
		Server.SetLevel(slog.LevelFatal)
	default:
		Server.SetLevel(slog.LevelInfo)
	}

	switch strings.ToLower(serverConfig.Log.Output) {
	case "stdout":
		Server.SetOutputToStdout()
	case "stderr":
		Server.SetOutputToStderr()
	case "file":
		name := serverConfig.Log.File.Name
		extensionName := serverConfig.Log.File.ExtensionName
		addDate := serverConfig.Log.File.AddDate

		Server.SetOutputToFile(name, extensionName, addDate)
	}

	Server.SetWithCallerInfo(serverConfig.Log.WithCallerInfo)
}
