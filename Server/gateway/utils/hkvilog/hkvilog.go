package hkvilog

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// 日志级别
const (
	DEBUG = iota
	INFO
	WARN
	ERROR
	FATAL
)

// 日志级别名称
var levelNames = []string{
	"DEBUG",
	"INFO",
	"WARN",
	"ERROR",
	"FATAL",
}

// 日志级别颜色（ANSI转义码）
var levelColors = []string{
	"\033[36m", // DEBUG - 青色
	"\033[32m", // INFO - 绿色
	"\033[33m", // WARN - 黄色
	"\033[31m", // ERROR - 红色
	"\033[35m", // FATAL - 紫色
}

const colorReset = "\033[0m"

// Logger 日志记录器
type Logger struct {
	level  int
	logger *log.Logger
}

// 全局日志记录器
var globalLogger = &Logger{
	level:  INFO,
	logger: log.New(os.Stdout, "", 0),
}

// SetLevel 设置日志级别
func SetLevel(level int) {
	globalLogger.level = level
}

// logf 格式化输出日志
func (l *Logger) logf(level int, format string, v ...interface{}) {
	if level < l.level {
		return
	}

	// 获取调用者信息
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// 提取文件名
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' || file[i] == '\\' {
			file = file[i+1:]
			break
		}
	}

	// 格式化时间
	now := time.Now()
	timestamp := now.Format("2006-01-02 15:04:05")

	// 构建日志消息
	msg := fmt.Sprintf(format, v...)
	levelName := levelNames[level]
	levelColor := levelColors[level]

	// 输出带颜色的日志
	logLine := fmt.Sprintf("%s[%s]%s %s %s:%d %s",
		levelColor, levelName, colorReset,
		timestamp, file, line, msg)

	l.logger.Print(logLine)

	// 如果是FATAL级别，退出程序
	if level == FATAL {
		os.Exit(1)
	}
}

// Debug 调试日志
func Debug(v ...interface{}) {
	globalLogger.logf(DEBUG, fmt.Sprint(v...))
}

// Debugf 格式化调试日志
func Debugf(format string, v ...interface{}) {
	globalLogger.logf(DEBUG, format, v...)
}

// Info 信息日志
func Info(v ...interface{}) {
	globalLogger.logf(INFO, fmt.Sprint(v...))
}

// Infof 格式化信息日志
func Infof(format string, v ...interface{}) {
	globalLogger.logf(INFO, format, v...)
}

// Warn 警告日志
func Warn(v ...interface{}) {
	globalLogger.logf(WARN, fmt.Sprint(v...))
}

// Warnf 格式化警告日志
func Warnf(format string, v ...interface{}) {
	globalLogger.logf(WARN, format, v...)
}

// Error 错误日志
func Error(v ...interface{}) {
	globalLogger.logf(ERROR, fmt.Sprint(v...))
}

// Errorf 格式化错误日志
func Errorf(format string, v ...interface{}) {
	globalLogger.logf(ERROR, format, v...)
}

// Fatal 致命错误日志
func Fatal(v ...interface{}) {
	globalLogger.logf(FATAL, fmt.Sprint(v...))
}

// Fatalf 格式化致命错误日志
func Fatalf(format string, v ...interface{}) {
	globalLogger.logf(FATAL, format, v...)
}
