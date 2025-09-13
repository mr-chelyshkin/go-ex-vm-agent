package worker

import (
	"context"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusStopping  TaskStatus = "stopping"
	TaskStatusStopped   TaskStatus = "stopped"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusCompleted TaskStatus = "completed"
)

type WorkerStatus string

const (
	WorkerStatusIdle     WorkerStatus = "idle"
	WorkerStatusStarting WorkerStatus = "starting"
	WorkerStatusRunning  WorkerStatus = "running"
	WorkerStatusStopping WorkerStatus = "stopping"
	WorkerStatusStopped  WorkerStatus = "stopped"
	WorkerStatusFailed   WorkerStatus = "failed"
)

type Task interface {
	Name() string
	Run(ctx context.Context) error
	Stop(ctx context.Context) error
}

type TaskInfo struct {
	Name      string
	Status    TaskStatus
	StartedAt *time.Time
	StoppedAt *time.Time
	Error     error
}

type taskWrapper struct {
	task   Task
	info   *TaskInfo
	cancel context.CancelFunc
	done   chan struct{}
}
