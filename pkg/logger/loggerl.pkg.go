package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go-ecommerce-backend-api.com/pkg/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config setting.LoggerSetting) *LoggerZap {
	loglevel := config.Log_Level
		// debug, info, warn, error, dpanic, panic, fatal
	var level zapcore.Level
	switch loglevel {
	case "debug":
		{
			level = zapcore.DebugLevel
		}
	case "info":
		{
			level = zapcore.InfoLevel
		}
	case "warn":
		{
			level = zapcore.WarnLevel
		}
	case "error":
		{
			level = zapcore.ErrorLevel
		}
	case "dpanic":
		{
			level = zapcore.DPanicLevel
		}
	case "panic":
		{
			level = zapcore.PanicLevel
		}
	case "fatal":
		{
			level = zapcore.FatalLevel
		}
	default:
		{
			level = zapcore.InfoLevel
		}
	}
	encoder := getEncoderLog()
	//curent time
	currentDate := time.Now().Format("2006-01-02")
	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s_%s.log", config.File_log_path, config.File_log_name, currentDate), // log file name
		MaxSize:    config.Max_size,                                                                      // megabytes
		MaxBackups: config.Max_backups,
		MaxAge:     config.Max_age,  //days
		Compress:   config.Compress, // disabled by default
	}
	core := zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), level)

	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}
}
func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	// 1748664226.0961385 -> 2025-05-31T11:03:46.096+0700
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// ts-> time
	encodeConfig.TimeKey = "time"
	// form infor
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// "caller": cli/main.go:20
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)
}
