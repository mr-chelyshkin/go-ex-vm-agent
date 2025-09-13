package config

import "go-ex-vm-agent/internal/logger"

type loggerConfig struct {
	Level   logger.LogLevel  `mapstructure:"level"`
	Format  logger.LogFormat `mapstructure:"format"`
	Output  logger.LogOutput `mapstructure:"output"`
	Options loggerOptions    `mapstructure:"options"`
}

type loggerOptions struct {
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxSize    int    `mapstructure:"max_size"`
	Compress   bool   `mapstructure:"compress"`
	FilePath   string `mapstructure:"path"`
}

func (lc loggerConfig) ToPkgConfig() logger.Config {
	return logger.Config{
		Level:      lc.Level,
		Format:     lc.Format,
		Output:     lc.Output,
		Path:       lc.Options.FilePath,
		MaxBackups: lc.Options.MaxBackups,
		MaxAge:     lc.Options.MaxAge,
		MaxSize:    lc.Options.MaxSize,
		Compress:   lc.Options.Compress,
	}
}
