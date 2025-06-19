package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/natefinch/lumberjack"
	"github.com/olekukonko/tablewriter"
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

// PrintStartupBanner displays a professional startup banner using tablewriter
func (l *LoggerZap) PrintStartupBanner(serviceName, version, port, env string) {
	// Print header with colors
	headerColor := color.New(color.FgCyan, color.Bold)
	headerColor.Println("\n🛒 E-COMMERCE BACKEND API")
	headerColor.Println(strings.Repeat("═", 50))

	// Create table for server information
	table := tablewriter.NewWriter(os.Stdout)

	// Set headers and data
	table.Header("PROPERTY", "VALUE")
	table.Append("🏷️ Service Name", serviceName)
	table.Append("📦 Version", version)
	table.Append("🌍 Environment", strings.ToUpper(env))
	table.Append("🚪 Port", port)
	table.Append("⏰ Started At", time.Now().Format("2006-01-02 15:04:05"))

	// Render the table
	table.Render()

	// Print startup message
	startupColor := color.New(color.FgGreen, color.Bold)
	startupColor.Println("\n🚀 Starting E-Commerce Backend API server...")
	fmt.Println(strings.Repeat("─", 50))
}

// LogInitStep logs từng bước khởi tạo với format đẹp
func (l *LoggerZap) LogInitStep(component string, success bool, err error) {
	if success {
		successColor := color.New(color.FgGreen, color.Bold)
		successColor.Printf("✅ %-30s ", component)
		color.New(color.FgGreen).Println("initialized successfully")
		l.Info(fmt.Sprintf("%s initialized successfully", component))
	} else {
		errorColor := color.New(color.FgRed, color.Bold)
		errorColor.Printf("❌ %-30s ", component)
		color.New(color.FgRed).Println("initialization failed")
		l.Error(fmt.Sprintf("Failed to initialize %s", component), zap.Error(err))
	}
}

// LogInitStart logs bắt đầu khởi tạo một component với format đẹp
func (l *LoggerZap) LogInitStart(component string) {
	initColor := color.New(color.FgYellow, color.Bold)
	initColor.Printf("🔧 Initializing %-25s", component)
	color.New(color.FgYellow).Println("...")
	l.Info(fmt.Sprintf("Initializing %s...", component))
}

// LogDBConnection logs kết nối database với table format
func (l *LoggerZap) LogDBConnection(dbType string, host string, success bool, err error) {
	if success {
		// Create a small table for successful DB connection
		table := tablewriter.NewWriter(os.Stdout)
		table.Header("DATABASE CONNECTION", "STATUS")
		table.Append(fmt.Sprintf("💾 %s", dbType), "✅ CONNECTED")
		table.Append("Host", host)
		table.Append("Time", time.Now().Format("15:04:05"))
		table.Render()

		l.Info(fmt.Sprintf("Connected to %s database", dbType),
			zap.String("host", host))
	} else {
		errorColor := color.New(color.FgRed, color.Bold)
		errorColor.Printf("💾 %-20s ", fmt.Sprintf("%s Database", strings.ToUpper(dbType)))
		color.New(color.FgRed).Println("❌ CONNECTION FAILED")
		color.New(color.FgRed).Printf("   Host: %s\n", host)
		if err != nil {
			color.New(color.FgRed).Printf("   Error: %s\n", err.Error())
		}

		l.Error(fmt.Sprintf("Failed to connect to %s database", dbType),
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

// LogShutdown logs khi server shutdown với format đẹp
func (l *LoggerZap) LogShutdown(reason string) {
	l.Info("Server shutting down gracefully...", zap.String("reason", reason))

	// Create shutdown table
	shutdownColor := color.New(color.FgYellow, color.Bold)
	shutdownColor.Println("\n🛑 SERVER SHUTDOWN")
	shutdownColor.Println(strings.Repeat("═", 40))

	table := tablewriter.NewWriter(os.Stdout)
	table.Header("SHUTDOWN INFO", "VALUE")
	table.Append("🛑 Status", "SHUTTING DOWN")
	table.Append("📝 Reason", reason)
	table.Append("⏰ Time", time.Now().Format("2006-01-02 15:04:05"))
	table.Render()

	goodbyeColor := color.New(color.FgCyan, color.Bold)
	goodbyeColor.Println("\n🛑 Server shutdown complete. Goodbye!")
	fmt.Println(strings.Repeat("═", 40))
}
