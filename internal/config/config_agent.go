package config

import (
	"go-ex-vm-agent/internal/runner"
	"go-ex-vm-agent/internal/worker"
	"time"
)

type agentConfig struct {
	RestartOptions agentRestartOptions `mapstructure:"restart_options"`
	TaskOptions    agentTaskOptions    `mapstructure:"task_options"`
	WorkersTimeout time.Duration       `mapstructure:"graceful_shutdown_workers_timeout"`
	RunnerTimeout  time.Duration       `mapstructure:"graceful_shutdown_agent_timeout"`
}

type agentRestartOptions struct {
	Delay            time.Duration `mapstructure:"delay"`
	MaxRestarts      int           `mapstructure:"max_restarts"`
	RestartExponent  bool          `mapstructure:"restart_exponent"`
	RestartOnFailure bool          `mapstructure:"restart_on_failure"`
}

type agentTaskOptions struct {
	MaxTimeout    time.Duration `mapstructure:"max_task_timeout"`
	MaxCount      int           `mapstructure:"max_task_count"`
	StopOnFailure bool          `mapstructure:"stop_on_failure"`
}

func (ac agentConfig) ToRunnerConfig() runner.Config {
	return runner.Config{
		EnableRestart:      ac.RestartOptions.RestartOnFailure,
		ExponentialBackoff: ac.RestartOptions.RestartExponent,
		MaxRestarts:        ac.RestartOptions.MaxRestarts,
		RestartDelay:       ac.RestartOptions.Delay,
		ShutdownTimeout:    ac.RunnerTimeout,
	}
}

func (ac agentConfig) ToWorkerConfig() worker.Config {
	return worker.Config{
		StopOnError:     ac.TaskOptions.StopOnFailure,
		TaskStopTimeout: ac.TaskOptions.MaxTimeout,
		TaskTimeout:     ac.TaskOptions.MaxTimeout,
		MaxTasks:        ac.TaskOptions.MaxCount,
	}
}
