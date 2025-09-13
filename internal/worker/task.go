package worker

import (
	"context"
	"time"
)

// BaseTask предоставляет базовую реализацию Task интерфейса
type BaseTask struct {
	name string
}

func NewBaseTask(name string) *BaseTask {
	return &BaseTask{name: name}
}

func (t *BaseTask) Name() string {
	return t.name
}

func (t *BaseTask) Run(ctx context.Context) error {
	// Базовая реализация - ничего не делает
	<-ctx.Done()
	return ctx.Err()
}

func (t *BaseTask) Stop(ctx context.Context) error {
	// Базовая реализация - просто возвращает nil
	return nil
}

// TickerTask - задача, которая выполняется по таймеру
type TickerTask struct {
	*BaseTask
	interval time.Duration
	handler  func(ctx context.Context) error
}

func NewTickerTask(name string, interval time.Duration, handler func(ctx context.Context) error) *TickerTask {
	return &TickerTask{
		BaseTask: NewBaseTask(name),
		interval: interval,
		handler:  handler,
	}
}

func (t *TickerTask) Run(ctx context.Context) error {
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := t.handler(ctx); err != nil {
				return executionError("ticker task '%s' failed: %v", t.Name(), err)
			}
		}
	}
}

// OnceTask - задача, которая выполняется один раз
type OnceTask struct {
	*BaseTask
	handler func(ctx context.Context) error
}

func NewOnceTask(name string, handler func(ctx context.Context) error) *OnceTask {
	return &OnceTask{
		BaseTask: NewBaseTask(name),
		handler:  handler,
	}
}

func (t *OnceTask) Run(ctx context.Context) error {
	return t.handler(ctx)
}
