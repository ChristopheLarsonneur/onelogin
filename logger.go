package onelogin

import (
	"github.com/op/go-logging"
	"os"
)

// LoggerModule is onelogin
const LoggerModule = "onelogin"

// Logger
var logger = logging.MustGetLogger(LoggerModule)

func init() {
	SetLogLevel(logging.WARNING)
}

// SetLogLevel define the level of log
func SetLogLevel(level logging.Level) {
	logging.SetFormatter(LogFormat())
	backend := logging.AddModuleLevel(StderrBackend(""))
	backend.SetLevel(level, "")
	logger.SetBackend(backend)
}

// LogFormat setup
func LogFormat() logging.Formatter {
	return logging.MustStringFormatter(
		`%{level:.8s} %{shortfile} %{shortfunc} ▶  %{message}`,
	)
}

// StderrBackend to display backend errors
func StderrBackend(prefix string) logging.Backend {
	return logging.NewLogBackend(os.Stderr, prefix, 0)
}

// StderrFormatter is the log formatter
func StderrFormatter(prefix string) {
	logging.SetFormatter(LogFormat())
}
