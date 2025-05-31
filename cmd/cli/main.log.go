package main

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// // zap logger
	// sugar := zap.NewExample().Sugar()
	// sugar.Infof("hello name: %s, age: %d", "thientriel", 18)
	// // logger

	// // logger := zap.NewExample()
	// // logger.Info("hello world", zap.String("name", "TipGo"), zap.Int("age", 18))

	// // logger development
	// logger, _ := zap.NewDevelopment()
	// logger.Info("hello Development logger")

	// // logger production
	// logger, _ = zap.NewProduction()
	// logger.Info("hello Production logger")
	encoder := getEncoderLog()
	sync := getWriterSync()
	core := zapcore.NewCore(encoder, sync, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	logger.Info("Info Log", zap.Int("line", 1))
	logger.Error("Error Log", zap.Int("line", 2))
}

// format logger
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

// writer logger in console and file
func getWriterSync() zapcore.WriteSyncer {
	file, _ := os.OpenFile("./log/log.txt", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	syncFile := zapcore.AddSync(file)
	syncConsole := zapcore.AddSync(os.Stderr)
	return zapcore.NewMultiWriteSyncer(syncConsole, syncFile)
}
