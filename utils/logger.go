package utils

import (
	"os"
	"path"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// var lock = &sync.Mutex{}
var Logger *zap.SugaredLogger

func InitLogger(env string) {
	// writerSyncer := getLogWriter()

	var core zapcore.Core

	if env == "development" {
		core = zapcore.NewTee(
			// zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(getJsonEncoder(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			// zapcore.NewCore(getFileEncoder(), writerSyncer, zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewTee(
			// zapcore.NewCore(getConsoleEncoder(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			zapcore.NewCore(getJsonEncoder(), zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
			// zapcore.NewCore(getFileEncoder(), writerSyncer, zapcore.DebugLevel),
		)
	}

	logger := zap.New(core, zap.AddCaller())

	Logger = logger.Sugar()
}

func getConsoleEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getJsonEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.StacktraceKey = "stacktrace"

	return zapcore.NewJSONEncoder(encoderConfig)
}

func getFileEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	logFilePath := viper.GetString("log.file.path")
	logFileName := viper.GetString("log.file.name")
	logFileMaxSize := viper.GetInt("log.file.maxsize")
	logFileMaxBackups := viper.GetInt("log.file.maxbackup")
	logFileMaxAge := viper.GetInt("log.file.maxage")
	logFile := path.Join(logFilePath, logFileName)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    logFileMaxSize,
		MaxBackups: logFileMaxBackups,
		MaxAge:     logFileMaxAge,
		Compress:   true,
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
