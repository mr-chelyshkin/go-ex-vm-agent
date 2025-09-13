package runner

import "fmt"

const (
	// ErrRunnerInit represents an error message format for failures during the initialization of a runner.
	ErrRunnerInit = "failed to initialize runner: %s"

	// ErrRunnerStart represents an error message format for failures occurring during the start of a runner.
	ErrRunnerStart = "failed to start runner: %s"

	// ErrRunnerStop represents an error message indicating a failure to stop the runner, formatted with additional context.
	ErrRunnerStop = "failed to stop runner: %s"

	// ErrRunnerRestart is an error message format for failures occurring during the restart of a runner.
	ErrRunnerRestart = "failed to restart runner: %s"

	// ErrWorkerManage represents a formatted error string for worker management-related failures.
	ErrWorkerManage = "worker management error: %s"

	// ErrSignalHandle indicates an error occurred during signal handling, formatted with an additional descriptive message.
	ErrSignalHandle = "signal handling error: %s"
)

func initError(format string, args ...any) error {
	return fmt.Errorf(ErrRunnerInit, fmt.Sprintf(format, args...))
}

func startError(format string, args ...any) error {
	return fmt.Errorf(ErrRunnerStart, fmt.Sprintf(format, args...))
}

func stopError(format string, args ...any) error {
	return fmt.Errorf(ErrRunnerStop, fmt.Sprintf(format, args...))
}

func restartError(format string, args ...any) error {
	return fmt.Errorf(ErrRunnerRestart, fmt.Sprintf(format, args...))
}

func workerManageError(format string, args ...any) error {
	return fmt.Errorf(ErrWorkerManage, fmt.Sprintf(format, args...))
}

func signalHandleError(format string, args ...any) error {
	return fmt.Errorf(ErrSignalHandle, fmt.Sprintf(format, args...))
}
