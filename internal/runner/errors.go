package runner

import "fmt"

const (
	ErrRunnerInit    = "failed to initialize runner: %s"
	ErrRunnerStart   = "failed to start runner: %s"
	ErrRunnerStop    = "failed to stop runner: %s"
	ErrRunnerRestart = "failed to restart runner: %s"
	ErrWorkerManage  = "worker management error: %s"
	ErrSignalHandle  = "signal handling error: %s"
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
