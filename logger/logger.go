package logger

import "fmt"
import "os"
import "runtime"
import "path"
import "strings"

type Logger struct {
	name  string
	level int
}

const (
	NONE = iota
	FATAL
	ERROR
	INFO
	DEBUG
	TRACE
)

var TTYGreen string
var TTYRed string
var TTYReset string

func init() {
	if os.ExpandEnv("$TERM") == "xterm" {
		TTYGreen = "\x1b[32m"
		TTYRed = "\x1b[31m"
		TTYReset = "\x1b[0m"
	}
}

func NewLogger(name string, level int) *Logger {
	/*if len(name) > 6 {
		name = name[0:6]
	}*/
	return &Logger{name: name, level: level}
}

func LevelFromString(level string) int {
	switch strings.ToUpper(level) {
	case "FATAL":
		return FATAL
	case "ERROR":
		return ERROR
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "TRACE":
		return TRACE
	}
	return -1
}

func levelString(level int) string {
	switch level {
	case FATAL:
		return "FATAL"
	case ERROR:
		return "ERROR"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return " INFO"
	case TRACE:
		return "TRACE"
	}
	return ""
}

func (log *Logger) print(dev *os.File, level int, msg string) {
	if log.level >= level {
		//fmt.Fprintf(dev, "%s: %6s: %s", levelString(level), log.name, msg)
		fmt.Fprintf(dev, "%s%s", log.name, msg)
		//fmt.Fprintf(dev, "%s", msg)
	}
}

func (log *Logger) Fatal(msg ...interface{}) {
	buf := fmt.Sprintln(msg...)
	log.print(os.Stderr, FATAL, TTYRed+buf+TTYReset)
	os.Exit(1)
}

func (log *Logger) Error(msg ...interface{}) {
	buf := fmt.Sprintln(msg...)
	log.print(os.Stderr, ERROR, TTYRed+buf+TTYReset)
}

func (log *Logger) Info(msg ...interface{}) {
	buf := fmt.Sprintln(msg...)
	log.print(os.Stdout, INFO, buf)
}

func (log *Logger) Message(msg ...interface{}) {
	buf := fmt.Sprintln(msg...)
	log.print(os.Stdout, INFO, TTYGreen+buf+TTYReset)
}

func (log *Logger) Debug(msg ...interface{}) {
	buf := fmt.Sprintln(msg...)
	log.print(os.Stdout, DEBUG, buf)
}

func (log *Logger) Trace(msg ...interface{}) {
	buf := fmt.Sprintln(msg...)

	pc, file, line, _ := runtime.Caller(1)
	fnc := runtime.FuncForPC(pc)
	stacktrace := fmt.Sprintf("%s:%d:%s", path.Base(file), line, fnc.Name())

	buf = fmt.Sprintf("=== TRACE: %s === \n%s===============\n", stacktrace, buf)
	log.print(os.Stdout, TRACE, buf)
}

func (log *Logger) Debugf(format string, msg ...interface{}) {
	buf := fmt.Sprintf(format, msg...)
	log.print(os.Stdout, DEBUG, buf)
}
