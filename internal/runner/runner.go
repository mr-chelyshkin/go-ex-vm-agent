package runner

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go-ex-vm-agent/internal/logger"
	"go-ex-vm-agent/internal/worker"
)

type Runner struct {
	config       Config
	workerConfig worker.Config
	logger       *logger.Logger
	taskFactory  TaskFactory

	mu           sync.RWMutex
	status       RunnerStatus
	worker       *worker.Worker
	restartCount int
	lastError    error

	ctx    context.Context
	cancel context.CancelFunc

	signals *signalHandler
	doneCh  chan struct{}
}

func New(config Config, workerConfig worker.Config, logger *logger.Logger, taskFactory TaskFactory) (*Runner, error) {
	if err := config.Validate(); err != nil {
		return nil, initError(err.Error())
	}

	if err := workerConfig.Validate(); err != nil {
		return nil, initError("worker config validation failed: %v", err)
	}

	if logger == nil {
		return nil, initError("logger cannot be nil")
	}

	if taskFactory == nil {
		return nil, initError("task factory cannot be nil")
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Runner{
		config:       config,
		workerConfig: workerConfig,
		logger:       logger,
		taskFactory:  taskFactory,
		status:       RunnerStatusIdle,
		ctx:          ctx,
		cancel:       cancel,
		signals:      newSignalHandler(),
		doneCh:       make(chan struct{}),
	}, nil
}

// Start запускает runner и начинает обработку сигналов
func (r *Runner) Start() error {
	r.mu.Lock()
	if r.status != RunnerStatusIdle {
		r.mu.Unlock()
		return startError("runner is not idle, current status: %s", r.status)
	}
	r.status = RunnerStatusStarting
	r.mu.Unlock()

	r.logger.Info().Msg("Starting runner")

	// Устанавливаем обработку сигналов
	if err := r.setupSignalHandling(); err != nil {
		r.mu.Lock()
		r.status = RunnerStatusFailed
		r.lastError = err
		r.mu.Unlock()
		return startError("failed to setup signal handling: %v", err)
	}

	// Запускаем worker
	if err := r.startWorker(); err != nil {
		r.mu.Lock()
		r.status = RunnerStatusFailed
		r.lastError = err
		r.mu.Unlock()
		return startError("failed to start worker: %v", err)
	}

	r.mu.Lock()
	r.status = RunnerStatusRunning
	r.mu.Unlock()

	r.logger.Info().Msg("Runner started successfully")

	// Запускаем основной цикл обработки
	go r.run()

	return nil
}

// Stop выполняет graceful shutdown runner'а
func (r *Runner) Stop() error {
	r.mu.Lock()
	if r.status != RunnerStatusRunning && r.status != RunnerStatusRestarting {
		currentStatus := r.status
		r.mu.Unlock()
		return stopError("runner is not running, current status: %s", currentStatus)
	}
	r.status = RunnerStatusStopping
	r.mu.Unlock()

	r.logger.Info().Msg("Stopping runner")

	// Отменяем основной контекст
	r.cancel()

	// Сигнализируем shutdown
	select {
	case r.signals.shutdown <- struct{}{}:
	default:
	}

	// Ждем завершения
	<-r.doneCh

	r.mu.Lock()
	r.status = RunnerStatusStopped
	r.mu.Unlock()

	r.logger.Info().Msg("Runner stopped successfully")
	return nil
}

// Restart перезапускает runner
func (r *Runner) Restart() error {
	r.logger.Info().Msg("Restarting runner")

	select {
	case r.signals.restart <- struct{}{}:
		return nil
	default:
		return restartError("restart already in progress")
	}
}

// GetInfo возвращает информацию о состоянии runner'а
func (r *Runner) GetInfo() RunnerInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	info := RunnerInfo{
		Status:       r.status,
		RestartCount: r.restartCount,
		LastError:    r.lastError,
	}

	if r.worker != nil {
		info.WorkerStatus = r.worker.GetStatus()
		info.WorkerTasks = r.worker.GetTasksInfo()
	}

	return info
}

// Wait блокируется до завершения runner'а
func (r *Runner) Wait() {
	<-r.doneCh
}

// setupSignalHandling настраивает обработку системных сигналов
func (r *Runner) setupSignalHandling() error {
	sigChan := make(chan os.Signal, 1)

	// SIGINT, SIGTERM - graceful shutdown
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// SIGUSR1 - restart (если поддерживается системой)
	if supportsSignal(syscall.SIGUSR1) {
		signal.Notify(sigChan, syscall.SIGUSR1)
	}

	// SIGHUP - reload config (если поддерживается системой)
	if supportsSignal(syscall.SIGHUP) {
		signal.Notify(sigChan, syscall.SIGHUP)
	}

	go func() {
		for {
			select {
			case <-r.ctx.Done():
				signal.Stop(sigChan)
				return
			case sig := <-sigChan:
				r.handleSystemSignal(sig)
			}
		}
	}()

	return nil
}

// handleSystemSignal обрабатывает системные сигналы
func (r *Runner) handleSystemSignal(sig os.Signal) {
	r.logger.Info().
		Str("signal", sig.String()).
		Msg("Received system signal")

	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		select {
		case r.signals.shutdown <- struct{}{}:
		default:
		}
	case syscall.SIGUSR1:
		select {
		case r.signals.restart <- struct{}{}:
		default:
		}
	case syscall.SIGHUP:
		select {
		case r.signals.reload <- struct{}{}:
		default:
		}
	}
}

// run основной цикл runner'а
func (r *Runner) run() {
	defer close(r.doneCh)

	for {
		select {
		case <-r.ctx.Done():
			r.shutdownWorker()
			return

		case <-r.signals.shutdown:
			r.shutdownWorker()
			return

		case <-r.signals.restart:
			if err := r.restartWorker(); err != nil {
				r.logger.Error().
					Err(err).
					Msg("Failed to restart worker")

				if !r.config.EnableRestart || r.shouldStopRestarting() {
					r.shutdownWorker()
					return
				}
			}

		case <-r.signals.reload:
			r.logger.Info().Msg("Config reload requested")
			// TODO: реализовать перезагрузку конфигурации
		}
	}
}

// startWorker запускает worker с задачами
func (r *Runner) startWorker() error {
	w, err := worker.New(r.workerConfig)
	if err != nil {
		return workerManageError("failed to create worker: %v", err)
	}

	// Регистрируем задачи
	tasks := r.taskFactory()
	for _, task := range tasks {
		if err := w.RegisterTask(task); err != nil {
			return workerManageError("failed to register task '%s': %v", task.Name(), err)
		}
	}

	// Запускаем worker
	if err := w.Start(r.ctx); err != nil {
		return workerManageError("failed to start worker: %v", err)
	}

	r.mu.Lock()
	r.worker = w
	r.mu.Unlock()

	// Мониторим состояние worker'а
	go r.monitorWorker()

	return nil
}

// shutdownWorker останавливает worker
func (r *Runner) shutdownWorker() {
	r.mu.RLock()
	w := r.worker
	r.mu.RUnlock()

	if w == nil {
		return
	}

	r.logger.Info().Msg("Shutting down worker")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), r.config.ShutdownTimeout)
	defer cancel()

	if err := w.Stop(shutdownCtx); err != nil {
		r.logger.Error().
			Err(err).
			Msg("Failed to shutdown worker gracefully")
	}

	r.mu.Lock()
	r.worker = nil
	r.mu.Unlock()
}

// restartWorker перезапускает worker
func (r *Runner) restartWorker() error {
	r.mu.Lock()
	r.status = RunnerStatusRestarting
	r.restartCount++
	r.mu.Unlock()

	r.logger.Info().
		Int("attempt", r.restartCount).
		Msg("Restarting worker")

	// Останавливаем текущий worker
	r.shutdownWorker()

	// Ждем задержку перед рестартом
	if r.config.RestartDelay > 0 {
		time.Sleep(r.config.RestartDelay)
	}

	// Запускаем новый worker
	if err := r.startWorker(); err != nil {
		r.mu.Lock()
		r.lastError = err
		r.mu.Unlock()
		return err
	}

	r.mu.Lock()
	r.status = RunnerStatusRunning
	r.lastError = nil
	r.mu.Unlock()

	r.logger.Info().Msg("Worker restarted successfully")
	return nil
}

// monitorWorker отслеживает состояние worker'а
func (r *Runner) monitorWorker() {
	r.mu.RLock()
	w := r.worker
	r.mu.RUnlock()

	if w == nil {
		return
	}

	w.Wait()

	// Worker завершился - проверяем нужен ли рестарт
	if r.config.EnableRestart && !r.shouldStopRestarting() {
		r.logger.Warn().Msg("Worker stopped unexpectedly, initiating restart")

		select {
		case r.signals.restart <- struct{}{}:
		case <-r.ctx.Done():
		}
	}
}

// shouldStopRestarting проверяет, следует ли прекратить попытки рестарта
func (r *Runner) shouldStopRestarting() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.config.MaxRestarts > 0 && r.restartCount >= r.config.MaxRestarts
}

// supportsSignal проверяет поддержку сигнала системой
func supportsSignal(sig syscall.Signal) bool {
	// Простая проверка - на Windows не все UNIX сигналы поддерживаются
	switch sig {
	case syscall.SIGUSR1, syscall.SIGHUP:
		// Эти сигналы не поддерживаются на Windows
		return os.Getenv("GOOS") != "windows"
	default:
		return true
	}
}
