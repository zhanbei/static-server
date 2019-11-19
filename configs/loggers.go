package configs

import (
	"os"
)

type LoggerFormat string

const (
	// Combined Text
	LoggerFormatText LoggerFormat = "text"
	// JSON with Rich Data
	LoggerFormatJson LoggerFormat = "json"
)

type LoggerTarget = string

const (
	LoggerTargetStdout LoggerTarget = "stdout"

	LoggerTargetFile LoggerTarget = "${file}"

	LoggerTargetDir LoggerTarget = "${dir}"
)

type OptionLoggers = []*OptionLogger

// Logger
type OptionLogger struct {
	Enabled bool `json:"enabled"`
	// text | json
	Format LoggerFormat `json:"format"`
	// A log file per host.
	// You are recommended to turn it off, if the hosts received kinds of arbitrary or dynamic,
	// which means that infinite number of logging files may be created.
	// FIX-ME Group similar hosts into a single file may be supported.
	PerHost bool `json:"perHost"`
	// Whether to print the content to be logged to stdout.
	Stdout bool `json:"stdout"`
	// The value should be `stdout|${file}` if the `perHost` option is false.
	// The value should be `${dir}` if the `perHost` option is true.
	Target LoggerTarget `json:"target"`
	// The recorder instance to be used to record requests logs.
	//Recorder IRecorder `json:"-"`

	LogWriter *os.File `json:"-"`
}

func NewLogger(format LoggerFormat, stdout bool, target string) *OptionLogger {
	return &OptionLogger{true, format, false, stdout, target, nil}
}

func (m *OptionLogger) IsValid() bool {
	if m.Format != LoggerFormatText || m.Format != LoggerFormatJson {
		return false
	}
	if m.Target != LoggerTargetStdout {
		return false
	}
	return true
}

type OptionLoggerGorilla struct {
	Enabled bool `json:"enabled"`
	// "combined"
	Format string `json:"format"`

	Stdout bool `json:"stdout"`

	Target LoggerTarget `json:"target"`

	LogWriter *os.File `json:"-"`
}