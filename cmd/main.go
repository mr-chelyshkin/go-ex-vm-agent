package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	agent "go-ex-vm-agent"
	"go-ex-vm-agent/internal/config"
	"go-ex-vm-agent/internal/logger"
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

	// Создаем основной контекст приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Настраиваем graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Создаем и настраиваем worker
	workerConfig := worker.Config{
		ShutdownTimeout: 30 * time.Second,
		TaskTimeout:     5 * time.Minute,
		MaxTasks:        10,
		StopOnError:     false,
	}

	w, err := worker.New(workerConfig)
	if err != nil {
		agent.Logger.Fatal().Err(err).Msg("Failed to create worker")
	}

	// Регистрируем задачи
	configWatcher := worker.NewTickerTask("config-watcher", 5*time.Second, func(ctx context.Context) error {
		agent.Logger.Debug().Msg("Config watcher tick")
		// TODO: логика отслеживания конфига
		return nil
	})

	healthCheck := worker.NewTickerTask("health-check", 30*time.Second, func(ctx context.Context) error {
		agent.Logger.Debug().Msg("Health check tick")
		// TODO: логика health check
		return nil
	})

	if err := w.RegisterTask(configWatcher); err != nil {
		agent.Logger.Fatal().Err(err).Msg("Failed to register config watcher task")
	}

	if err := w.RegisterTask(healthCheck); err != nil {
		agent.Logger.Fatal().Err(err).Msg("Failed to register health check task")
	}

	// Запускаем worker
	if err := w.Start(ctx); err != nil {
		agent.Logger.Fatal().Err(err).Msg("Failed to start worker")
	}

	agent.Logger.Info().Msg("Application started successfully")

	// Ждем сигнал завершения
	<-sigChan
	agent.Logger.Info().Msg("Shutdown signal received")

	// Отменяем основной контекст
	cancel()

	// Создаем контекст с таймаутом для graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Останавливаем worker
	if err := w.Stop(shutdownCtx); err != nil {
		agent.Logger.Error().Err(err).Msg("Failed to stop worker gracefully")
		os.Exit(1)
	}

	agent.Logger.Info().Msg("Application stopped gracefully")
}
