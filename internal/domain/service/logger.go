package service

type Logger interface {
	Info(message string, args ...any)
	Error(message string, args ...any)
}
