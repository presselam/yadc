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
	LabelTrace          = "TRACE"
	LabelDebug          = "DEBUG"
	LabelInfo           = "INFO "
	LabelWarn           = "WARN "
	LabelError          = "ERROR"
)

var logLevel LogLevel = LevelInfo
var logColor bool = true

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

func logMessage(level string, color string, args ...any) {
	if !logColor {
		color = ansi.DefaultFG
	}

	var message string
	if len(args) > 0 {
		values := []string{}
		for _, v := range args {
			values = append(values, fmt.Sprintf("%v", v))
		}
		message = strings.Join(values, "")
	}

	log.Printf("%s%-5s - %s%s", color, level, message, ansi.Reset)
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

	logMessage(LabelTrace, ansi.Magenta, caller)
}

func Debug(args ...any) {
	if logLevel > LevelDebug {
		return
	}
	logMessage(LabelDebug, ansi.Yellow, args)
}

func Info(args ...any) {
	if logLevel > LevelInfo {
		return
	}
	logMessage(LabelInfo, ansi.DefaultFG, args)
}

func Warn(args ...any) {
	if logLevel > LevelWarn {
		return
	}
	logMessage(LabelWarn, ansi.LightRed, args)
}

func Error(args ...any) {
	if logLevel > LevelError {
		return
	}
	logMessage(LabelError, ansi.Red, args)
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
