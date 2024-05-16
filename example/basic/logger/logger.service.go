package logger

import (
	"fmt"
	"io"
	"time"
)

type LOG_LEVEL int
type LOG_LEVEL_STRING string

const localTimeFormat = "2006-01-02 15:04:05"

// https://docs.spring.io/spring-boot/docs/2.1.13.RELEASE/reference/html/boot-features-logging.html
const (
	LOG_LEVEL_TRACE LOG_LEVEL = iota
	LOG_LEVEL_DEBUG
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
	LOG_LEVEL_FATAL
	LOG_LEVEL_PANIC
)

var logLevelStrings = []string{
	"[Trace]",
	"[Debug]",
	"[Info ]",
	"[Warn ]",
	"[Error]",
	"[Fatal]",
	"[Panic]",
}

var allowLogLevel = LOG_LEVEL_TRACE

func (lv LOG_LEVEL) ToString() string {
	if lv >= LOG_LEVEL_TRACE && lv <= LOG_LEVEL_PANIC {
		return logLevelStrings[lv]
	}
	return fmt.Sprintf("[?%d] ", lv)
}

type logPayload struct {
	Module  *string
	Message *string
	Type    LOG_LEVEL
}

var isInit = false

var stdOut *io.Writer
var stdErr *io.Writer

func buildLog(logType LOG_LEVEL, module *string, message ...any) {
	if logType < allowLogLevel {
		return
	}

	str := fmt.Sprintf(
		"%s %s [%s] %s",
		time.Now().Format(localTimeFormat),
		logType.ToString(),
		*module,
		fmt.Sprint(message...),
	)

	var writer *io.Writer
	switch logType {
	case LOG_LEVEL_TRACE, LOG_LEVEL_DEBUG, LOG_LEVEL_INFO:
		writer = stdOut
		break
	case LOG_LEVEL_WARN, LOG_LEVEL_ERROR, LOG_LEVEL_FATAL:
		writer = stdErr
		break
	case LOG_LEVEL_PANIC:
		writer = stdErr
		fmt.Fprintln(*writer, str)
		panic(str)
	default:
		writer = stdOut
	}

	fmt.Fprintln(*writer, str)
}

func InitLogger(
	out *io.Writer,
	err *io.Writer,
) {
	if isInit {
		return
	}

	isInit = true
	stdOut = out
	stdErr = err
}

func SetLevel(l LOG_LEVEL) {
	allowLogLevel = l
}

type Logger struct {
	module *string
}

func NewLogger(module string) *Logger {
	l := new(Logger)
	l.module = &module

	return l
}

func (l *Logger) ErrorF(format string, params ...any) {
	l.Error(fmt.Sprintf(format, params...))
}

func (l *Logger) Error(message ...any) {
	buildLog(LOG_LEVEL_ERROR, l.module, message...)
}

func (l *Logger) Warn(message ...any) {
	buildLog(LOG_LEVEL_WARN, l.module, message...)
}

func (l *Logger) WarnF(format string, params ...any) {
	l.Warn(fmt.Sprintf(format, params...))
}

func (l *Logger) Info(message ...any) {
	buildLog(LOG_LEVEL_INFO, l.module, message...)
}

func (l *Logger) InfoF(format string, params ...any) {
	l.Info(fmt.Sprintf(format, params...))
}

func (l *Logger) Debug(message ...any) {
	buildLog(LOG_LEVEL_DEBUG, l.module, message...)
}

func (l *Logger) DebugF(format string, params ...any) {
	l.Debug(fmt.Sprintf(format, params...))
}

func (l *Logger) Trace(message ...any) {
	buildLog(LOG_LEVEL_TRACE, l.module, message...)
}

func (l *Logger) TraceF(format string, params ...any) {
	l.Trace(fmt.Sprintf(format, params...))
}

func (l *Logger) Fatal(message ...any) {
	buildLog(LOG_LEVEL_FATAL, l.module, message...)
}

func (l *Logger) FatalF(format string, params ...any) {
	l.Fatal(fmt.Sprintf(format, params...))
}

func (l *Logger) Panic(message ...any) {
	buildLog(LOG_LEVEL_PANIC, l.module, message...)
}

func (l *Logger) PanicF(format string, params ...any) {
	l.Trace(fmt.Sprintf(format, params...))
}
