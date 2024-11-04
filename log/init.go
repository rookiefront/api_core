package log

import (
	"github.com/rookiefront/api-core/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var writeLevel zapcore.Level = -2
var printLevel zapcore.Level = -2
var printLogger *zap.Logger

// 用于记录详细的调试信息，帮助开发人员了解程序的运行细节。
// 适合在开发环境和测试环境中使用，但通常不会在生产环境中开启，以避免日志量过大。
var debugLogger *zap.Logger

// 用于记录程序的常规操作信息，如启动、停止、配置信息等。
// 适合用于了解程序的运行情况，一般会在生产环境开启。
var infoLogger *zap.Logger

// 用于记录可能会引发问题的事件，但不会立即影响程序运行的警告信息。
// 适合监控程序中的潜在问题，帮助开发人员提前预防。
var warnLogger *zap.Logger

// 用于记录程序在运行中遇到的错误，例如无法连接数据库、文件未找到等。
// 一般在生产环境中会开启，以便在发生问题时能及时了解并解决
var errorLogger *zap.Logger

// 表示“Development Panic”，是一个特别的日志级别，仅在开发环境中会触发 panic，用于记录一些严重的错误。
// 在生产环境中，它的行为类似于 Error，不会触发 panic，但在开发环境中可以立即中断程序执行，帮助快速定位问题。
var dPanicLogger *zap.Logger

// 表示严重错误，记录完日志后会触发 panic，导致程序中断并产生堆栈跟踪信息。
// 一般用于记录无法继续运行的严重问题
var panicLogger *zap.Logger

// 表示致命错误，记录完日志后会立即调用 os.Exit(1) 退出程序。
// 适合在遇到不可恢复的错误时使用，如配置文件加载失败等
var fatalLogger *zap.Logger

func SetMinWriteLevel(level zapcore.Level) {
	writeLevel = level
}
func SetMinPrintLevel(level zapcore.Level) {
	printLevel = level
}
func init() {
	printLogger = newLoggerWithColor()
}
func Debug(msg string, fields ...zap.Field) {
	level := zapcore.DebugLevel
	if debugLogger == nil {
		debugLogger = newLogger("logs/debug.log", level)
		defer debugLogger.Sync()
	}
	if printLevel <= level {
		printLogger.Debug(msg, fields...)
	}
	if writeLevel <= level {
		debugLogger.Debug(msg, fields...)
	}
}
func Info(msg string, fields ...zap.Field) {
	level := zapcore.InfoLevel
	if infoLogger == nil {
		infoLogger = newLogger("logs/info.log", level)
		defer infoLogger.Sync()
	}
	if printLevel <= level {
		printLogger.Info(msg, fields...)
	}
	if writeLevel <= level {
		infoLogger.Info(msg, fields...)
	}
}

func Warn(msg string, fields ...zap.Field) {
	level := zapcore.WarnLevel
	if warnLogger == nil {
		warnLogger = newLogger("logs/info.log", level)
		defer warnLogger.Sync()
	}
	if printLevel <= level {
		printLogger.Warn(msg, fields...)
	}
	if writeLevel <= level {
		warnLogger.Warn(msg, fields...)
	}
}

func Err(msg string, fields ...zap.Field) {
	level := zapcore.ErrorLevel
	if errorLogger == nil {
		errorLogger = newLogger("logs/error.log", level)
		defer errorLogger.Sync()
	}
	if printLevel <= level {
		printLogger.Error(msg, fields...)
	}
	if writeLevel <= level {
		errorLogger.Error(msg, fields...)
	}
}

func DPanic(msg string, fields ...zap.Field) {
	if !config.IsDev() {
		return
	}
	level := zapcore.DPanicLevel
	if dPanicLogger == nil {
		dPanicLogger = newLogger("logs/dev-panic.log", level)
		defer dPanicLogger.Sync()
	}
	if printLevel <= level {
		printLogger.DPanic(msg, fields...)
	}
	if writeLevel <= level {
		dPanicLogger.DPanic(msg, fields...)
	}
}
func Panic(msg string, fields ...zap.Field) {
	level := zapcore.PanicLevel
	if panicLogger == nil {
		panicLogger = newLogger("logs/panic.log", level)
		defer panicLogger.Sync()
	}
	if printLevel <= level {
		printLogger.Panic(msg, fields...)
	}
	if writeLevel <= level {
		panicLogger.Panic(msg, fields...)
	}
}
func Fatal(msg string, fields ...zap.Field) {
	level := zapcore.FatalLevel
	if fatalLogger == nil {
		fatalLogger = newLogger("logs/fatal.log", level)
		defer fatalLogger.Sync()
	}
	fatalLogger.Error(msg, fields...)
	if printLevel <= level {
		printLogger.Panic(msg, fields...)
	}
	if writeLevel <= level {
		fatalLogger.Panic(msg, fields...)
	}
}

func newLogger(filename string, level zapcore.Level) *zap.Logger {
	// 使用 lumberjack 进行日志轮转
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename, // 日志文件路径
		MaxSize:    100,      // 单个日志文件的最大尺寸 (MB)
		MaxAge:     30,       // 保留的最大天数
		MaxBackups: 7,        // 保留的旧文件数
		LocalTime:  true,
		Compress:   true, // 压缩旧文件
	})

	// 配置 zap 的编码器和级别
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "T",
		LevelKey:      "L",
		NameKey:       "N",
		CallerKey:     "C",
		MessageKey:    "M",
		StacktraceKey: "S",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 使用 JSON 编码器
		writer,                                // 使用 lumberjack 进行日志输出
		level,                                 // 日志级别
	)

	return zap.New(core, zap.AddCaller())
}

func newLoggerWithColor() *zap.Logger {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "N",
		CallerKey:      "C",
		MessageKey:     "M",
		StacktraceKey:  "S",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 使用颜色编码
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.DebugLevel)

	return zap.New(core, zap.AddCaller())
}
