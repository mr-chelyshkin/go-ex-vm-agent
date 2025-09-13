package main

import (
	agent "go-ex-vm-agent"
	"go-ex-vm-agent/internal/config"
	"go-ex-vm-agent/internal/logger"
)

func main() {
	cfg, err := config.Load("config.yaml")
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
