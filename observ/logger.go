package observ

import "go.uber.org/zap"

// TODO: (Logger) Close() error => zLogger.Sync()

type Logger struct {
	zLogger *zap.Logger
}

func NewLogger(config zap.Config) (*Logger, error) {
	zLogger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{zLogger}, nil
}

func (l *Logger) Named(name string) *Logger        { return &Logger{l.zLogger.Named(name)} }
func (l *Logger) With(fields ...zap.Field) *Logger { return &Logger{l.zLogger.With(fields...)} }
func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	return &Logger{l.zLogger.WithOptions(opts...)}
}

func (l *Logger) Error(message string, fields ...zap.Field) { l.zLogger.Error(message, fields...) }
func (l *Logger) Warn(message string, fields ...zap.Field)  { l.zLogger.Warn(message, fields...) }
func (l *Logger) Info(message string, fields ...zap.Field)  { l.zLogger.Info(message, fields...) }
func (l *Logger) Debug(message string, fields ...zap.Field) { l.zLogger.Debug(message, fields...) }

func (l *Logger) Sync() error { return l.zLogger.Sync() }
