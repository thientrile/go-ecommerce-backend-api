package logger

import (
	"fmt"
	"os"

	"github.com/natefinch/lumberjack"
	"go-ecommerce-backend-api.com/pkg/setting"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger(config setting.LoggerSetting, server setting.ServerSetting) *LoggerZap {
	loglevel := config.Log_Level
	version := server.Version
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

	// Tạo encoder cho file (JSON format)
	fileEncoder := getFileEncoderLog()
	// Tạo encoder cho console (human-readable format)
	consoleEncoder := getConsoleEncoderLog()

	hook := lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s_%s.log", config.File_log_path, config.File_log_name, version), // log file name
		MaxSize:    config.Max_size,                                                                  // megabytes
		MaxBackups: config.Max_backups,
		MaxAge:     config.Max_age,  //days
		Compress:   config.Compress, // disabled by default
	}

	// Console core với format đẹp
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level)
	// File core với JSON format
	fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(&hook), level)

	// Kết hợp cả 2 cores
	core := zapcore.NewTee(consoleCore, fileCore)

	return &LoggerZap{zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))}
}

// getFileEncoderLog tạo encoder JSON cho file logs
func getFileEncoderLog() zapcore.Encoder {
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

// getConsoleEncoderLog tạo encoder đẹp cho console
func getConsoleEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewDevelopmentEncoderConfig()

	// Thời gian format đẹp: 15:04:05.000
	encodeConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05.000")
	encodeConfig.TimeKey = "T"

	// Level với màu sắc
	encodeConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		var levelStr string
		switch level {
		case zapcore.DebugLevel:
			levelStr = "\033[36m[DBG]\033[0m" // Cyan
		case zapcore.InfoLevel:
			levelStr = "\033[32m[INF]\033[0m" // Green
		case zapcore.WarnLevel:
			levelStr = "\033[33m[WRN]\033[0m" // Yellow
		case zapcore.ErrorLevel:
			levelStr = "\033[31m[ERR]\033[0m" // Red
		case zapcore.DPanicLevel:
			levelStr = "\033[35m[DPF]\033[0m" // Magenta
		case zapcore.PanicLevel:
			levelStr = "\033[35m[PNC]\033[0m" // Magenta
		case zapcore.FatalLevel:
			levelStr = "\033[31m[FTL]\033[0m" // Red Bold
		default:
			levelStr = fmt.Sprintf("[%s]", level.CapitalString())
		}
		enc.AppendString(levelStr)
	}

	// Caller format ngắn gọn
	encodeConfig.EncodeCaller = func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString("\033[90m" + caller.TrimmedPath() + "\033[0m") // Gray
	}

	// Message format
	encodeConfig.ConsoleSeparator = " "

	return zapcore.NewConsoleEncoder(encodeConfig)
}

// Helper methods cho việc logging dễ dàng hơn

// PrintStartupBanner in ra banner khi khởi động server
func (l *LoggerZap) PrintStartupBanner(serviceName, version, port, env string) {
	banner := fmt.Sprintf(`
╔══════════════════════════════════════════════════════════════╗
║                🛒 E-COMMERCE BACKEND API                     ║
║══════════════════════════════════════════════════════════════║
║  Service: %-20s                               ║
║  Version: %-20s                               ║
║  Environment: %-15s                               ║
║  Port: %-23s                               ║
╚══════════════════════════════════════════════════════════════╝`,
		serviceName, version, env, port)

	fmt.Println("\033[32m" + banner + "\033[0m") // Green color
	l.Info("🚀 Starting E-Commerce Backend API server...")
}

// LogInitStep logs từng bước khởi tạo
func (l *LoggerZap) LogInitStep(component string, success bool, err error) {
	if success {
		l.Info(fmt.Sprintf("✅ %s initialized successfully", component))
	} else {
		l.Error(fmt.Sprintf("❌ Failed to initialize %s", component), zap.Error(err))
	}
}

// LogInitStart logs bắt đầu khởi tạo một component
func (l *LoggerZap) LogInitStart(component string) {
	l.Info(fmt.Sprintf("🔧 Initializing %s...", component))
}

// LogDBConnection logs kết nối database
func (l *LoggerZap) LogDBConnection(dbType string, host string, success bool, err error) {
	if success {
		l.Info(fmt.Sprintf("💾 Connected to %s database", dbType),
			zap.String("host", host))
	} else {
		l.Error(fmt.Sprintf("💾 Failed to connect to %s database", dbType),
			zap.String("host", host),
			zap.Error(err))
	}
}

// LogHTTPRequest logs HTTP requests với format đẹp
func (l *LoggerZap) LogHTTPRequest(method, path, ip string, status int, duration string) {
	var statusIcon string
	switch {
	case status >= 200 && status < 300:
		statusIcon = "✅"
	case status >= 300 && status < 400:
		statusIcon = "↩️"
	case status >= 400 && status < 500:
		statusIcon = "⚠️"
	default:
		statusIcon = "❌"
	}

	l.Info(fmt.Sprintf("%s %s %s", statusIcon, method, path),
		zap.String("ip", ip),
		zap.Int("status", status),
		zap.String("duration", duration))
}

// LogShutdown logs khi server shutdown
func (l *LoggerZap) LogShutdown(reason string) {
	l.Info("🛑 Server shutting down gracefully...", zap.String("reason", reason))
	fmt.Println("\033[33m══════════════════════════════════════════════════════════════\033[0m")
	fmt.Println("\033[33m🛑 Server shutdown complete. Goodbye!\033[0m")
}
