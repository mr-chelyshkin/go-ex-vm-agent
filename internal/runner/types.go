package runner

import (
	"go-ex-vm-agent/internal/worker"
)

// RunnerStatus представляет статус runner'а
type RunnerStatus string

const (
	RunnerStatusIdle       RunnerStatus = "idle"
	RunnerStatusStarting   RunnerStatus = "starting"
	RunnerStatusRunning    RunnerStatus = "running"
	RunnerStatusStopping   RunnerStatus = "stopping"
	RunnerStatusStopped    RunnerStatus = "stopped"
	RunnerStatusRestarting RunnerStatus = "restarting"
	RunnerStatusFailed     RunnerStatus = "failed"
)

// TaskFactory функция для создания задач
type TaskFactory func() []worker.Task

// RunnerInfo содержит информацию о состоянии runner'а
type RunnerInfo struct {
	Status       RunnerStatus
	RestartCount int
	WorkerStatus worker.WorkerStatus
	WorkerTasks  map[string]worker.TaskInfo
	LastError    error
}

// SignalAction тип действия при получении сигнала
type SignalAction string

const (
	SignalActionShutdown SignalAction = "shutdown"
	SignalActionRestart  SignalAction = "restart"
	SignalActionReload   SignalAction = "reload"
)

// signalHandler внутренняя структура для обработки сигналов
type signalHandler struct {
	shutdown chan struct{}
	restart  chan struct{}
	reload   chan struct{}
}

func newSignalHandler() *signalHandler {
	return &signalHandler{
		shutdown: make(chan struct{}, 1),
		restart:  make(chan struct{}, 1),
		reload:   make(chan struct{}, 1),
	}
}
