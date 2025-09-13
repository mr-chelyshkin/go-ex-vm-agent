package worker

import "fmt"

const (
	ErrWorkerInit       = "failed to initialize worker: %s"
	ErrWorkerStart      = "failed to start worker: %s"
	ErrWorkerStop       = "failed to stop worker: %s"
	ErrTaskRegistration = "failed to register task: %s"
	ErrTaskExecution    = "task execution error: %s"
	ErrTaskTimeout      = "task timeout: %s"
)

func initError(format string, args ...any) error {
	return fmt.Errorf(ErrWorkerInit, fmt.Sprintf(format, args...))
}

func startError(format string, args ...any) error {
	return fmt.Errorf(ErrWorkerStart, fmt.Sprintf(format, args...))
}

func stopError(format string, args ...any) error {
	return fmt.Errorf(ErrWorkerStop, fmt.Sprintf(format, args...))
}

func registrationError(format string, args ...any) error {
	return fmt.Errorf(ErrTaskRegistration, fmt.Sprintf(format, args...))
}

func executionError(format string, args ...any) error {
	return fmt.Errorf(ErrTaskExecution, fmt.Sprintf(format, args...))
}

func timeoutError(format string, args ...any) error {
	return fmt.Errorf(ErrTaskTimeout, fmt.Sprintf(format, args...))
}
