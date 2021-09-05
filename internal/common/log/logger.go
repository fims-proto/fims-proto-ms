package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/spf13/viper"
)

const (
	debugLv = iota
	infoLv
	errLv
)

var (
	debug *log.Logger
	info  *log.Logger
	err   *log.Logger

	logEnabler func(lvl int, v ...interface{}) bool

	flagDebug = viper.GetBool("debug")
)

func NewStdLoggerAdapter() io.Writer {
	return os.Stderr
}

func NewStdLogEnablerAdapter() func(lvl int, v ...interface{}) bool {
	return func(lvl int, _ ...interface{}) bool {
		return lvl > 0 || flagDebug
	}
}

func InitLoggers(logEnablerAdapter func(lvl int, v ...interface{}) bool, loggerAdapters ...io.Writer) {
	out := io.MultiWriter(loggerAdapters...)

	const (
		sd = "[ debug ]: "
		si = "[ info ] : "
		se = "[ error ]: "
	)

	debug = log.New(out, sd, log.LstdFlags|log.Lshortfile)
	info = log.New(out, si, log.LstdFlags)
	err = log.New(out, se, log.LstdFlags)

	logEnabler = logEnablerAdapter

	Infoln("logger initiated")

	if flagDebug {
		Debugln("debug mode enabled")
	}
}

func Debugln(v ...interface{}) {
	logln(debug, debugLv, v...)
}

func Debugf(format string, v ...interface{}) {
	logf(debug, debugLv, format, v...)
}

func Infoln(v ...interface{}) {
	logln(info, infoLv, v...)
}

func Infof(format string, v ...interface{}) {
	logf(info, infoLv, format, v...)
}

func Errorln(v ...interface{}) {
	logln(err, errLv, v...)
}

func Errorf(format string, v ...interface{}) {
	logf(err, errLv, format, v...)
}

func logln(logger *log.Logger, lvl int, args ...interface{}) {
	if !logEnabler(lvl) {
		return
	}
	_ = logger.Output(3, fmt.Sprintln(args...))
}

func logf(logger *log.Logger, lvl int, format string, args ...interface{}) {
	if !logEnabler(lvl) {
		return
	}
	_ = logger.Output(3, fmt.Sprintf(format, args...))
}
