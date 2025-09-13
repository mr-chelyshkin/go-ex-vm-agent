package logger

import "strings"

// LogLevel represents the severity level of a log message.
type LogLevel string

// IsValid checks if the LogLevel is one of the predefined valid log levels.
func (l LogLevel) IsValid() bool {
	for _, valid := range ValidLogLevels {
		if l == valid {
			return true
		}
	}
	return false
}

// LogFormat represents the format of log output.
type LogFormat string

// IsValid checks if the LogFormat is one of the predefined valid log formats in the validLogFormats list.
func (f LogFormat) IsValid() bool {
	for _, valid := range validLogFormats {
		if f == valid {
			return true
		}
	}
	return false
}

// LogOutput represents the destination where log messages are written.
type LogOutput string

// IsValid checks if the LogOutput is one of the predefined valid log outputs in the validLogOutputs list.
func (o LogOutput) IsValid() bool {
	for _, valid := range validLogOutputs {
		if o == valid {
			return true
		}
	}
	return false
}

const (
	// LevelDebug represents the debug severity level in logging, used for detailed informational events helpful during development.
	LevelDebug LogLevel = "debug"

	// LevelInfo represents the log level for general informational messages.
	LevelInfo LogLevel = "info"

	// LevelWarn represents a warning severity level in the logging system, indicating potential issues that require attention.
	LevelWarn LogLevel = "warn"

	// LevelError represents the log level for errors, used to categorize log messages indicating failures or issues.
	LevelError LogLevel = "error"

	// LevelFatal represents the fatal log level, typically used for critical issues causing immediate application termination.
	LevelFatal LogLevel = "fatal"

	// LevelPanic represents a log level indicating a critical condition that causes the application to panic immediately.
	LevelPanic LogLevel = "panic"

	// LevelDisabled represents a log level where logging is completely disabled.
	LevelDisabled LogLevel = "disabled"
)

const (
	// FormatJSON represents the JSON format for log output.
	FormatJSON LogFormat = "json"

	// FormatConsole specifies the console output format for logging.
	FormatConsole LogFormat = "console"
)

const (
	// OutputJournal denotes the log output destination as a system journal like journald.
	OutputJournal LogOutput = "journal"

	// OutputStdout represents the log output destination as standard output (stdout).
	OutputStdout LogOutput = "stdout"

	// OutputStderr represents log output directed to the standard error stream.
	OutputStderr LogOutput = "stderr"

	// OutputFile represents a logging output that writes log messages to a file.
	OutputFile LogOutput = "file"
)

// ValidLogLevels defines a list of predefined log levels considered valid for the logging system.
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
