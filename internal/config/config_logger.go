package config

import "go-ex-vm-agent/internal/logger"

// loggerConfig represents the configuration settings for the logger.
type loggerConfig struct {
	// Level specifies the logging level for the logger configuration.
	Level logger.LogLevel `mapstructure:"level"`

	// Format specifies the format of the logs (e.g., JSON, plain text).
	Format logger.LogFormat `mapstructure:"format"`

	// Output specifies the destination where log messages should be written (e.g., file, stdout, stderr).
	Output logger.LogOutput `mapstructure:"output"`

	// FileOptions defines file-specific configuration options for logger output, such as file path, size, age, and backups.
	FileOptions loggerFileOptions `mapstructure:"options"`
}

// loggerFileOptions defines configuration for file-based logging.
type loggerFileOptions struct {
	// MaxBackups specifies the maximum number of backup log files to retain.
	MaxBackups int `mapstructure:"max_backups"`

	// MaxAge specifies the maximum number of days logs will be retained before being automatically deleted.
	MaxAge int `mapstructure:"max_age"`

	// MaxSize specifies the maximum size (in MB) of the log file before it is rotated.
	MaxSize int `mapstructure:"max_size"`

	// Compress indicates whether old log files should be compressed using gzip.
	Compress bool `mapstructure:"compress"`

	// FilePath specifies the file path where the log file will be stored.
	FilePath string `mapstructure:"path"`
}

// ToLoggerConfig transforms a loggerConfig instance into the logger.Config structure used by the logger package.
func (lc loggerConfig) ToLoggerConfig() logger.Config {
	return logger.Config{
		Level:      lc.Level,
		Format:     lc.Format,
		Output:     lc.Output,
		MaxAge:     lc.FileOptions.MaxAge,
		MaxSize:    lc.FileOptions.MaxSize,
		Compress:   lc.FileOptions.Compress,
		Path:       lc.FileOptions.FilePath,
		MaxBackups: lc.FileOptions.MaxBackups,
	}
}
