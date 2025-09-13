package config

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Logger loggerConfig `mapstructure:"logger"`
}

func Load(path string) (*Config, error) {
	v := viper.New()

	ext, err := getConfigFormatByExt(filepath.Ext(path))
	if err != nil {
		return nil, err
	}

	v.SetConfigType(string(ext))
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &cfg, nil
}
