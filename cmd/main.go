package main

import (
	agent "go-ex-vm-agent"
	"go-ex-vm-agent/internal/config"
	"go-ex-vm-agent/internal/logger"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("config", "", "Path to config file")
	pflag.Parse()

	viper.SetEnvPrefix(agent.EnvPrefix)
	viper.AutomaticEnv()
	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		panic(err)
	}

	cfg, err := config.Load(viper.GetString("config"))
	if err != nil {
		panic(err)
	}

	loggerConfig := cfg.Logger.ToPkgConfig()
	agent.Logger, err = logger.New(loggerConfig)
	if err != nil {
		panic(err)
	}

	agent.Logger.Info().Msg("Hello World")
}
