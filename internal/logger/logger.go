package logger

import (
	"fmt"
	"github.com/mgutz/ansi"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type LogLevel int

const (
	LevelTrace LogLevel = 0
	LevelDebug LogLevel = 5
	LevelInfo  LogLevel = 10
	LevelWarn  LogLevel = 15
	LevelError LogLevel = 20
)

var logLevel LogLevel

func Setup(level ...LogLevel) {
	if len(level) == 0 {
		ev := os.Getenv("LOG_LEVEL")
		switch strings.ToLower(ev) {
		case "trace":
			logLevel = LevelTrace
		case "debug":
			logLevel = LevelDebug
		case "info":
			logLevel = LevelInfo
		case "warn":
			logLevel = LevelWarn
		case "error":
			logLevel = LevelError
		default:
			logLevel = LevelInfo
		}
	} else {
		logLevel = level[0]
	}
}

func Trace(args ...any) {
	if logLevel > LevelTrace {
		return
	}

	pcs := make([]uintptr, 1)
	runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs)

	frame, _ := frames.Next()
	caller := filepath.Base(frame.Function)

	if len(args) > 0 {
		values := []string{}
		for _, v := range args {
			values = append(values, fmt.Sprintf("%v", v))
		}
		caller += " - [" + strings.Join(values, "][") + "]"
	}

	log.Printf("%sTRACE - %s%s", ansi.Magenta, caller, ansi.Reset)
}

func Debug(args ...string) {
	if logLevel > LevelDebug {
		return
	}

	msg := strings.Join(args, "")
	log.Printf("%sDEBUG - %s%s", ansi.Yellow, msg, ansi.Reset)
}

func Info(args ...string) {
	if logLevel > LevelInfo {
		return
	}

	msg := strings.Join(args, "")
	log.Printf("INFO  - %s", msg)
}

func Warn(args ...string) {
	if logLevel > LevelWarn {
		return
	}

	msg := strings.Join(args, "")
	log.Printf("%sWARN  - %s%s", ansi.LightRed, msg, ansi.Reset)
}

func Error(args ...string) {
	if logLevel > LevelError {
		return
	}

	msg := strings.Join(args, "")
	log.Printf("%sERROR - %s%s", ansi.Red, msg, ansi.Reset)
}

func StartBanner() {
	var s string
	switch logLevel {
	case LevelTrace:
		s = "Trace"
	case LevelDebug:
		s = "Debug"
	case LevelInfo:
		s = "Info"
	case LevelWarn:
		s = "Warn"
	case LevelError:
		s = "Error"
	default:
		s = "Info"
	}

	log.Println("#=======================================================")
	log.Printf("# [%s] - [%v]", s, time.Now())
	log.Println("#=======================================================")
}
