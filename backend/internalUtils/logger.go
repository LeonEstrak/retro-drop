package internalUtils

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/LeonEstrak/retro-drop/backend/constants"
)

// logLevel type for defining log levels
type logLevel int

const (
	debugLevel logLevel = iota
	infoLevel
	warningLevel
	errorLevel
	fatalLevel
	panicLevel
)

var logLevelNames = map[logLevel]string{
	debugLevel:   "DEBUG",
	infoLevel:    "INFO",
	warningLevel: "WARN",
	errorLevel:   "ERROR",
	fatalLevel:   "FATAL",
	panicLevel:   "PANIC",
}

// Logger struct
type Logger struct {
	level      logLevel
	logger     *log.Logger
	mu         sync.Mutex
	timeFormat string
}

var (
	once   sync.Once
	logger *Logger
)

// GetLogger returns a global Logger instance.
//
// The Logger instance is created with the configured log level. The first call
// to GetLogger will create the Logger instance and subsequent calls will return
// the same instance.
func GetLogger() *Logger {
	once.Do(func() {
		logger = newLogger(constants.LOG_LEVEL)
	})
	return logger
}

// newLogger returns a new Logger instance with the given log level.
//
// The log level is parsed from the given string and must be one of the
// following: "debug", "info", "warning", "error", "fatal", or "panic".
// If the log level is invalid, the logger will default to "info".
//
// The logger uses the given log level and writes to os.Stdout with the
// prefix "[RETRO-DROP] ". The time format for the logger is fixed to
// "2006-01-02 15:04:05.000".
func newLogger(levelStr string) *Logger {
	level := parseLogLevel(levelStr)
	logger := log.New(os.Stdout, "[RETRO-DROP] ", 0) // Remove default flags, we'll format ourselves
	return &Logger{
		level:      level,
		logger:     logger,
		timeFormat: "2006-01-02 15:04:05.000",
	}
}

// parseLogLevel parses a string representing a log level and returns the corresponding logLevel constant.
//
// The input string is case-insensitive and may be one of the following:
// "debug", "info", "warn", "warning", "error", "fatal", or "panic".
// If the string is invalid, the function will print a warning and return the default log level, which is infoLevel.
func parseLogLevel(levelStr string) logLevel {
	levelStr = strings.ToLower(levelStr)
	switch levelStr {
	case "debug":
		return debugLevel
	case "info":
		return infoLevel
	case "warn", "warning":
		return warningLevel
	case "error":
		return errorLevel
	case "fatal":
		return fatalLevel
	case "panic":
		return panicLevel
	default:
		fmt.Printf("Warning: Invalid log level '%s', defaulting to INFO\n", levelStr)
		return infoLevel
	}
}

func (l *Logger) log(level logLevel, format string, v ...interface{}) {
	if level < l.level {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60)) // Bengaluru Time
	levelName := logLevelNames[level]
	message := fmt.Sprintf(format, v...)
	logOutput := fmt.Sprintf("%s [%s] %s", now.Format(l.timeFormat), levelName, message)
	l.logger.Println(logOutput)

	if level == panicLevel {
		panic(message)
	}
}

// Debug logs a message at Debug level.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.log(debugLevel, format, v...)
}

// Info logs a message at Info level.
func (l *Logger) Info(format string, v ...interface{}) {
	l.log(infoLevel, format, v...)
}

// Warn logs a message at Warning level.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.log(warningLevel, format, v...)
}

// Error logs a message at Error level.
func (l *Logger) Error(format string, v ...interface{}) {
	l.log(errorLevel, format, v...)
}

// Fatal logs a message at Fatal level and then calls os.Exit(1).
func (l *Logger) Fatal(format string, v ...interface{}) {
	l.log(fatalLevel, format, v...)
}

// Panic logs a message at Panic level and then calls panic().
func (l *Logger) Panic(format string, v ...interface{}) {
	l.log(panicLevel, format, v...)
}
