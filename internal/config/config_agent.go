package config

import (
	"time"

	"go-ex-vm-agent/internal/runner"
	"go-ex-vm-agent/internal/worker"
)

// agentConfig defines the configuration settings for an agent.
type agentConfig struct {
	// RestartOptions defines the restart configuration options for the agent.
	RestartOptions agentRestartOptions `mapstructure:"restart_options"`

	// TaskOptions defines settings for task execution.
	TaskOptions agentTaskOptions `mapstructure:"task_options"`

	// WorkersTimeout specifies the duration to wait for workers to complete during a graceful shutdown process.
	WorkersTimeout time.Duration `mapstructure:"graceful_shutdown_workers_timeout"`

	// RunnerTimeout specifies the duration allowed for the agent to shut down gracefully before being forcefully terminated.
	RunnerTimeout time.Duration `mapstructure:"graceful_shutdown_agent_timeout"`
}

// agentRestartOptions defines configuration settings related to restarting an agent.
type agentRestartOptions struct {
	// Delay specifies the time duration to wait before attempting an agent restart.
	Delay time.Duration `mapstructure:"delay"`

	// MaxRestarts specifies the maximum number of times the agent will attempt to restart.
	MaxRestarts int `mapstructure:"max_restarts"`

	// RestartExponent determines whether an exponential backoff delay strategy is applied between restart attempts.
	RestartExponent bool `mapstructure:"restart_exponent"`

	// RestartOnFailure determines whether the agent should automatically restart upon encountering a failure or crash.
	RestartOnFailure bool `mapstructure:"restart_on_failure"`
}

// agentTaskOptions defines configuration options for controlling task execution behavior in the agent.
type agentTaskOptions struct {
	// MaxTimeout specifies the maximum duration allowed for a task to execute before it is forcibly terminated.
	MaxTimeout time.Duration `mapstructure:"max_task_timeout"`

	// MaxCount specifies the maximum number of tasks that can be executed concurrently.
	MaxCount int `mapstructure:"max_task_count"`

	// StopOnFailure determines if task execution should stop when a failure is encountered.
	StopOnFailure bool `mapstructure:"stop_on_failure"`
}

// ToRunnerConfig transforms an agentConfig instance into the runner.Config structure used by the runner package.
func (ac agentConfig) ToRunnerConfig() runner.Config {
	return runner.Config{
		EnableRestart:      ac.RestartOptions.RestartOnFailure,
		ExponentialBackoff: ac.RestartOptions.RestartExponent,
		MaxRestarts:        ac.RestartOptions.MaxRestarts,
		RestartDelay:       ac.RestartOptions.Delay,
		ShutdownTimeout:    ac.RunnerTimeout,
	}
}

// ToWorkerConfig transforms an agentConfig instance into the worker.Config structure used by the worker package.
func (ac agentConfig) ToWorkerConfig() worker.Config {
	return worker.Config{
		StopOnError:     ac.TaskOptions.StopOnFailure,
		TaskStopTimeout: ac.TaskOptions.MaxTimeout,
		TaskTimeout:     ac.TaskOptions.MaxTimeout,
		MaxTasks:        ac.TaskOptions.MaxCount,
	}
}
