package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application's configuration.
type Config struct {
	// Logger represents the configuration for the logging system.
	Logger loggerConfig `mapstructure:"logger"`

	// Agent represents the configuration settings for the agent.
	Agent agentConfig `mapstructure:"agent"`
}

// Load reads and parses a configuration file from the specified path and returns a Config object or an error.
func Load(path string) (*Config, error) {
	ext, err := getConfigFormatByExt(filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(string(ext))
	if err = v.ReadInConfig(); err != nil {
		return nil, initError("read file error: %s", err.Error())
	}

	var cfg Config
	if err = v.Unmarshal(&cfg); err != nil {
		return nil, initError("serialization error: %s", err.Error())
	}
	return &cfg, nil
}
