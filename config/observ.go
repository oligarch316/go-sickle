package config

import (
	"fmt"

	"github.com/oligarch316/go-sickle/observ"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TODO: ObservDataLog as it's own type when/if necessary

type ObservData struct {
	LogEncoding string `dhall:"logEncoding"`
	LogLevel    string `dhall:"logLevel"`

	LogCaller     bool `dhall:"logCaller"`
	LogStacktrace bool `dhall:"logStacktrace"`
}

func MergeObservData(base, priority ObservData) ObservData {
	res := ObservData{LogLevel: base.LogLevel, LogEncoding: base.LogEncoding}

	if priority.LogLevel != "" {
		res.LogLevel = priority.LogLevel
	}

	if priority.LogEncoding != "" {
		res.LogEncoding = priority.LogEncoding
	}

	return res
}

func BuildObserv(data ObservData) (*observ.Logger, error) {
	atomicLevel, err := zap.ParseAtomicLevel(data.LogLevel)
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

	switch data.LogEncoding {
	case "console":
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	case "json":
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		encoderConfig.EncodeTime = zapcore.EpochTimeEncoder
		encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	default:
		return nil, fmt.Errorf("invalid encoding: %s", data.LogEncoding)
	}

	zConfig := zap.Config{
		Level:             atomicLevel,
		DisableCaller:     !data.LogCaller,
		DisableStacktrace: !data.LogStacktrace,
		Encoding:          data.LogEncoding,
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
	}

	return observ.NewLogger(zConfig)
}
