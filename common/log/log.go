package log

import (
	"strings"

	"github.com/base-server/go/common/config"
	"github.com/common-library/go/log/slog"
)

var Log slog.Log

func Initialize(kind string) {
	switch strings.ToLower(config.Get(kind + ".log.level").(string)) {
	case "trace":
		Log.SetLevel(slog.LevelTrace)
	case "debug":
		Log.SetLevel(slog.LevelDebug)
	case "info":
		Log.SetLevel(slog.LevelInfo)
	case "warn":
		Log.SetLevel(slog.LevelWarn)
	case "error":
		Log.SetLevel(slog.LevelError)
	case "fatal":
		Log.SetLevel(slog.LevelFatal)
	default:
		Log.SetLevel(slog.LevelInfo)
	}

	switch strings.ToLower(config.Get(kind + ".log.output").(string)) {
	case "stdout":
		Log.SetOutputToStdout()
	case "stderr":
		Log.SetOutputToStderr()
	case "file":
		name := config.Get(kind + ".log.file.name").(string)
		extensionName := config.Get(kind + ".log.file.extensionName").(string)
		addDate := config.Get(kind + ".log.file.addDate").(bool)

		Log.SetOutputToFile(name, extensionName, addDate)
	}

	Log.SetWithCallerInfo(config.Get(kind + ".log.withCallerInfo").(bool))
}
