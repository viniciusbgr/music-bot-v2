// TODO: Refactor the Logger to use a more structured approach, more flags and options
package logger

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

type TraceLevel uint8

const (
	TraceLevelInfo          TraceLevel = 1 << iota
	TraceLevelWarnings      TraceLevel = 1 << 2
	TraceLevelErrors        TraceLevel = 1 << 3
	TraceLevelDebugStandard TraceLevel = 1 << 4
	TraceLevelDebugAdvanced TraceLevel = 1 << 5

	// Standards flags
	StdTraceLevelStandard      TraceLevel = TraceLevelInfo | TraceLevelWarnings | TraceLevelErrors
	StdTraceLevelDebugOnlyData TraceLevel = TraceLevelDebugStandard | TraceLevelDebugAdvanced
	StdTraceLevelDebugAll      TraceLevel = TraceLevelDebugAdvanced | StdTraceLevelStandard

	ColorReset  string = "\033[0m"
	ColorRed    string = "\033[31m"
	ColorBlue   string = "\033[34m"
	ColorGreen  string = "\033[32m"
	ColorYellow string = "\033[33m"
)

type Logger struct {
	log       *log.Logger
	level     TraceLevel
	colorized bool
}

var (
	DefaultLogger *Logger

	ErrWriterNil = errors.New("logger: writer is nil")
)

func New(out io.Writer, level TraceLevel, colorized bool) (*Logger, error) {
	if out == nil {
		return nil, ErrWriterNil
	}

	if level == 0 {
		level = StdTraceLevelStandard
	}

	logger := log.New(out, "", log.LstdFlags)

	if level&TraceLevelDebugAdvanced == TraceLevelDebugAdvanced {
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}

	return &Logger{log: logger, level: level, colorized: colorized}, nil
}

func (l *Logger) Error(v ...any) error {
	if l.level&TraceLevelErrors != TraceLevelErrors {
		return nil
	}

	out := fmt.Append(nil, v...)

	return l.log.Output(2, ColorRed+" "+"[ERROR] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Errorf(format string, v ...any) error {
	if l.level&TraceLevelErrors != TraceLevelErrors {
		return nil
	}

	out := fmt.Append(nil, v...)

	return l.log.Output(2, ColorRed+" "+"[ERROR] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Warn(v ...any) error {
	if l.level&TraceLevelWarnings != TraceLevelWarnings {
		return nil
	}

	out := fmt.Append(nil, v...)
	return l.log.Output(2, ColorYellow+" "+"[WARN] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Warnf(format string, v ...any) error {
	if l.level&TraceLevelWarnings != TraceLevelWarnings {
		return nil
	}

	out := fmt.Append(nil, v...)

	return l.log.Output(2, ColorYellow+" "+"[WARN] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Info(v ...any) error {
	if l.level&TraceLevelInfo != TraceLevelInfo {
		return nil
	}

	out := fmt.Append(nil, v...)

	return l.log.Output(2, ColorGreen+" "+"[INFO] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Infof(format string, v ...any) error {
	if l.level&TraceLevelInfo != TraceLevelInfo {
		return nil
	}

	out := fmt.Appendf(nil, format, v...)

	return l.log.Output(2, ColorGreen+" "+"[INFO] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Debug(format string, v ...any) error {
	if l.level&TraceLevelDebugAdvanced != TraceLevelDebugAdvanced && l.level&TraceLevelDebugStandard != TraceLevelDebugStandard {
		return nil
	}

	out := fmt.Appendf(nil, format, v...)
	return l.log.Output(2, ColorBlue+" "+"[DEBUG] -> "+string(out)+ColorReset+"\n")
}

func (l *Logger) Fatal(v ...any) {
	out := fmt.Append(nil, v...)

	l.log.Output(2, ColorRed+" "+"[FATAL] -> "+string(out)+ColorReset+"\n")
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...any) {
	out := fmt.Appendf(nil, format, v...)

	l.log.Output(2, ColorRed+" "+"[FATAL] -> "+string(out)+ColorReset+"\n")
	os.Exit(1)
}
