package worker

import (
	"context"
	agent "go-ex-vm-agent"
	"sync"
	"time"
)

type Worker struct {
	tasks map[string]*taskWrapper

	// TODO: maybe atomic from pointers?
	mu     sync.RWMutex
	config Config
	status WorkerStatus

	stopCh chan struct{}
	doneCh chan struct{}
}

func New(config Config) (*Worker, error) {
	if err := config.Validate(); err != nil {
		return nil, initError(err.Error())
	}

	return &Worker{
		tasks: make(map[string]*taskWrapper),

		config: config,
		status: WorkerStatusIdle,
		stopCh: make(chan struct{}),
		doneCh: make(chan struct{}),
	}, nil
}

func (w *Worker) RegisterTask(task Task) error {
	if task == nil {
		return registrationError("task cannot be nil")
	}
	name := task.Name()
	if name == "" {
		return registrationError("task name cannot be empty")
	}
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.tasks[name]; exists {
		return registrationError("task '%s' already registered", name)
	}
	if len(w.tasks) >= w.config.MaxTasks {
		return registrationError("cannot register task '%s': maximum tasks limit (%d) reached", name, w.config.MaxTasks)
	}
	// TODO: maybe error status only
	if w.status != WorkerStatusIdle {
		return registrationError("cannot register task '%s': worker is not idle", name)
	}

	w.tasks[name] = &taskWrapper{
		task: task,
		info: &TaskInfo{
			Name:   name,
			Status: TaskStatusPending,
		},
		done: make(chan struct{}),
	}

	agent.Logger.Info().
		Str("task", name).
		Msg("Task registered successfully")
	return nil
}

func (w *Worker) Start(ctx context.Context) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// TODO: maybe list of available statuses??
	if w.status != WorkerStatusIdle {
		return startError("worker is not idle, current status: %s", w.status)
	}

	if len(w.tasks) == 0 {
		return startError("no tasks registered")
	}

	w.status = WorkerStatusStarting
	agent.Logger.Info().
		Int("task_count", len(w.tasks)).
		Msg("Starting worker")

	for name, wrapper := range w.tasks {
		if err := w.startTask(ctx, wrapper); err != nil {
			w.status = WorkerStatusFailed
			return startError("failed to start task '%s': %v", name, err)
		}
	}

	w.status = WorkerStatusRunning
	agent.Logger.Info().Msg("Worker started successfully")
	go w.monitor(ctx)
	return nil
}

func (w *Worker) Stop(ctx context.Context) error {
	w.mu.Lock()
	if w.status != WorkerStatusRunning {
		w.mu.Unlock()
		return stopError("worker is not running, current status: %s", w.status)
	}
	w.status = WorkerStatusStopping
	w.mu.Unlock()

	agent.Logger.Info().Msg("Stopping worker")
	shutdownCtx, cancel := context.WithTimeout(ctx, w.config.ShutdownTimeout)
	defer cancel()

	close(w.stopCh)
	//TODO: use new wg Format.
	var wg sync.WaitGroup
	w.mu.RLock()
	for name, wrapper := range w.tasks {
		wg.Add(1)
		go func(name string, wrapper *taskWrapper) {
			defer wg.Done()
			w.stopTask(shutdownCtx, wrapper)
		}(name, wrapper)
	}
	w.mu.RUnlock()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		agent.Logger.Info().Msg("All tasks stopped gracefully")
	case <-shutdownCtx.Done():
		agent.Logger.Warn().
			Dur("timeout", w.config.ShutdownTimeout).
			Msg("Graceful shutdown timeout exceeded")
	}

	w.mu.Lock()
	w.status = WorkerStatusStopped
	w.mu.Unlock()

	close(w.doneCh)
	agent.Logger.Info().Msg("Worker stopped")
	return nil
}

func (w *Worker) GetStatus() WorkerStatus {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.status
}

func (w *Worker) GetTasksInfo() map[string]TaskInfo {
	w.mu.RLock()
	defer w.mu.RUnlock()

	info := make(map[string]TaskInfo)
	for name, wrapper := range w.tasks {
		info[name] = *wrapper.info
	}
	return info
}

func (w *Worker) Wait() {
	<-w.doneCh
}

func (w *Worker) startTask(ctx context.Context, wrapper *taskWrapper) error {
	taskCtx, cancel := context.WithCancel(ctx)
	wrapper.cancel = cancel

	if w.config.TaskTimeout > 0 {
		taskCtx, cancel = context.WithTimeout(taskCtx, w.config.TaskTimeout)
		originalCancel := wrapper.cancel
		wrapper.cancel = func() {
			cancel()
			originalCancel()
		}
	}

	now := time.Now()
	wrapper.info.Status = TaskStatusRunning
	wrapper.info.StartedAt = &now

	go func() {
		defer close(wrapper.done)
		defer func() {
			now := time.Now()
			wrapper.info.StoppedAt = &now
		}()

		agent.Logger.Debug().
			Str("task", wrapper.task.Name()).
			Msg("Starting task")

		if err := wrapper.task.Run(taskCtx); err != nil {
			wrapper.info.Status = TaskStatusFailed
			wrapper.info.Error = err

			agent.Logger.Error().
				Err(err).
				Str("task", wrapper.task.Name()).
				Msg("Task failed")

			if w.config.StopOnError {
				go func() {
					if stopErr := w.Stop(context.Background()); stopErr != nil {
						agent.Logger.Error().
							Err(stopErr).
							Msg("Failed to stop worker after task error")
					}
				}()
			}
			return
		}

		wrapper.info.Status = TaskStatusCompleted
		agent.Logger.Debug().
			Str("task", wrapper.task.Name()).
			Msg("Task completed")
	}()
	return nil
}

func (w *Worker) stopTask(ctx context.Context, wrapper *taskWrapper) {
	wrapper.info.Status = TaskStatusStopping

	agent.Logger.Debug().
		Str("task", wrapper.task.Name()).
		Msg("Stopping task")

	if wrapper.cancel != nil {
		wrapper.cancel()
	}

	if err := wrapper.task.Stop(ctx); err != nil {
		agent.Logger.Warn().
			Err(err).
			Str("task", wrapper.task.Name()).
			Msg("Task stop returned error")
	}

	select {
	case <-wrapper.done:
		if wrapper.info.Status != TaskStatusFailed {
			wrapper.info.Status = TaskStatusStopped
		}
		agent.Logger.Debug().
			Str("task", wrapper.task.Name()).
			Msg("Task stopped")
	case <-ctx.Done():
		wrapper.info.Status = TaskStatusFailed
		wrapper.info.Error = timeoutError("task '%s' stop timeout", wrapper.task.Name())
		agent.Logger.Warn().
			Str("task", wrapper.task.Name()).
			Msg("Task stop timeout")
	}
}

func (w *Worker) monitor(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.stopCh:
			return
		case <-ticker.C:
			w.logTasksStatus()
		}
	}
}

func (w *Worker) logTasksStatus() {
	w.mu.RLock()
	defer w.mu.RUnlock()

	running := 0
	failed := 0
	completed := 0

	for _, wrapper := range w.tasks {
		switch wrapper.info.Status {
		case TaskStatusRunning:
			running++
		case TaskStatusFailed:
			failed++
		case TaskStatusCompleted:
			completed++
		}
	}

	agent.Logger.Debug().
		Int("total", len(w.tasks)).
		Int("running", running).
		Int("failed", failed).
		Int("completed", completed).
		Msg("Tasks status")
}
