package logger

import "strings"

type LogLevel string

func (l LogLevel) IsValid() bool {
	for _, valid := range ValidLogLevels {
		if l == valid {
			return true
		}
	}
	return false
}

type LogFormat string

func (f LogFormat) IsValid() bool {
	for _, valid := range validLogFormats {
		if f == valid {
			return true
		}
	}
	return false
}

type LogOutput string

func (o LogOutput) IsValid() bool {
	for _, valid := range validLogOutputs {
		if o == valid {
			return true
		}
	}
	return false
}

const (
	LevelDebug    LogLevel = "debug"
	LevelInfo     LogLevel = "info"
	LevelWarn     LogLevel = "warn"
	LevelError    LogLevel = "error"
	LevelFatal    LogLevel = "fatal"
	LevelPanic    LogLevel = "panic"
	LevelDisabled LogLevel = "disabled"
)

const (
	FormatJSON    LogFormat = "json"
	FormatConsole LogFormat = "console"
)

const (
	OutputJournal LogOutput = "journal"
	OutputStdout  LogOutput = "stdout"
	OutputStderr  LogOutput = "stderr"
	OutputFile    LogOutput = "file"
)

var ValidLogLevels = []LogLevel{
	LevelDisabled,
	LevelDebug,
	LevelInfo,
	LevelWarn,
	LevelError,
	LevelFatal,
	LevelPanic,
}

var validLogFormats = []LogFormat{
	FormatConsole,
	FormatJSON,
}

var validLogOutputs = []LogOutput{
	OutputStdout,
	OutputStderr,
	OutputFile,
}

func logLevelsString() string {
	levels := make([]string, len(ValidLogLevels))
	for i, level := range ValidLogLevels {
		levels[i] = string(level)
	}
	return strings.Join(levels, " ")
}

func logFormatsString() string {
	formats := make([]string, len(validLogFormats))
	for i, format := range validLogFormats {
		formats[i] = string(format)
	}
	return strings.Join(formats, " ")
}

func logOutputsString() string {
	outputs := make([]string, len(validLogOutputs))
	for i, output := range validLogOutputs {
		outputs[i] = string(output)
	}
	return strings.Join(outputs, " ")
}
