package logger

import (
	"log/slog"
)
type SlogLogger struct {
    logger *slog.Logger
}

func NewSlogLogger() *SlogLogger {
    return &SlogLogger{
        logger: slog.Default(),
    }
}

func (s *SlogLogger) Info(msg string, args ...any) {
    s.logger.Info(msg, args...)
}

func (s *SlogLogger) Debug(msg string, args ...any) {
    s.logger.Debug(msg, args...)
}

func (s *SlogLogger) Error(msg string, args ...any) {
    s.logger.Error(msg, args...)
}

func (s *SlogLogger) Warn(msg string, args ...any) {
    s.logger.Warn(msg, args...)
}


func (s *SlogLogger) Success(msg string, args ...any) {
    s.logger.Info(msg, args...)
}

func (s *SlogLogger) Fatal(msg string, args ...any) {
    s.logger.Error(msg, args...)
}

func (s *SlogLogger) Panic(msg string, args ...any) {
    s.logger.Error(msg, args...)
}


