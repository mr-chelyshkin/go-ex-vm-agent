package main

import (
	agent "go-ex-vm-agent"
	"go-ex-vm-agent/internal/logger"
)

func main() {
	lcfg := logger.Config{
		Level:      "",
		Format:     "console",
		Output:     "",
		Path:       "",
		MaxBackups: 0,
		MaxAge:     0,
		MaxSize:    0,
		Compress:   false,
	}
	var err error
	agent.Logger, err = logger.New(lcfg)

	if err != nil {
		panic(err)
	}

	agent.Logger.Info().Msg("Hello World")
}
