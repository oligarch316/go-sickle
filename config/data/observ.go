package data

import (
	"fmt"

	"github.com/oligarch316/go-sickle/config/value"
	"github.com/oligarch316/go-sickle/observ"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ObservLogConfig struct {
	Encoding value.String `dhall:"encoding"`
	Level    value.String `dhall:"level"`

	EnableCaller     value.Bool `dhall:"enableCaller"`
	EnableStacktrace value.Bool `dhall:"enableStacktrace"`
}

type ObservConfig struct {
	Log ObservLogConfig `dhall:"log"`
}

func BuildObserv(data ObservConfig) (*observ.Logger, error) {
	atomicLevel, err := zap.ParseAtomicLevel(string(data.Log.Level))
	if err != nil {
		return nil, err
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "M",
		StacktraceKey: "S",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	switch data.Log.Encoding {
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	case "json":
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
		encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	default:
		return nil, fmt.Errorf("invalid encoding: %s", string(data.Log.Encoding))
	}

	zConfig := zap.Config{
		Level:             atomicLevel,
		DisableCaller:     !bool(data.Log.EnableCaller),
		DisableStacktrace: !bool(data.Log.EnableStacktrace),
		Encoding:          string(data.Log.Encoding),
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return observ.NewLogger(zConfig)
}
