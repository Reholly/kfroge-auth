package implementation

import (
	"log/slog"
	"os"
	"sso-service/internal/domain/service"
)

const (
	FilePermissions = 0666
)

type Logger struct {
	slog.Logger
}

func NewLogger() (service.Logger, error) {
	/*file, err := os.OpenFile("app_logs", os.O_WRONLY|os.O_APPEND|os.O_CREATE, FilePermissions)
	if err != nil {
		return nil, err
	}
	*/
	return slog.New(slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		})), nil
}
