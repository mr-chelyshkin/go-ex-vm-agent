package main

import (
	"context"
	"time"

	agent "go-ex-vm-agent"
	"go-ex-vm-agent/internal/config"
	"go-ex-vm-agent/internal/logger"
	"go-ex-vm-agent/internal/runner"
	"go-ex-vm-agent/internal/worker"

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

	// Конфигурации для runner и worker
	runnerConfig := runner.Config{
		ShutdownTimeout: 30 * time.Second,
		RestartDelay:    5 * time.Second,
		MaxRestarts:     3,
		EnableRestart:   true,
	}

	workerConfig := worker.Config{
		ShutdownTimeout: 30 * time.Second,
		TaskTimeout:     5 * time.Minute,
		MaxTasks:        10,
		StopOnError:     false,
	}

	// Фабрика задач
	taskFactory := func() []worker.Task {
		return []worker.Task{
			worker.NewTickerTask("config-watcher", 5*time.Second, func(ctx context.Context) error {
				agent.Logger.Debug().Msg("Config watcher tick")
				// TODO: логика отслеживания конфига
				return nil
			}),
			worker.NewTickerTask("health-check", 30*time.Second, func(ctx context.Context) error {
				agent.Logger.Debug().Msg("Health check tick")
				// TODO: логика health check
				return nil
			}),
		}
	}

	// Создаем runner
	r, err := runner.New(runnerConfig, workerConfig, agent.Logger, taskFactory)
	if err != nil {
		agent.Logger.Fatal().Err(err).Msg("Failed to create runner")
	}

	// Запускаем runner
	if err := r.Start(); err != nil {
		agent.Logger.Fatal().Err(err).Msg("Failed to start runner")
	}

	agent.Logger.Info().Msg("Application started successfully")

	// Ждем завершения (runner сам обрабатывает сигналы)
	r.Wait()

	agent.Logger.Info().Msg("Application stopped gracefully")
}
