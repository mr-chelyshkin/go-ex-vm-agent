package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Logger loggerConfig `mapstructure:"logger"`
	Agent  agentConfig  `mapstructure:"agent"`
}

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
