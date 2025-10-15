package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel LogLevel
	logger       *log.Logger
)

// levelMap 用于将字符串转换为对应的日志级别
var levelMap = map[string]LogLevel{
	"debug": DEBUG,
	"info":  INFO,
	"warn":  WARN,
	"error": ERROR,
}

// colorMap 日志级别对应的终端颜色
var colorMap = map[LogLevel]string{
	DEBUG: "",         // 默认颜色
	INFO:  "\033[32m", // 绿色
	WARN:  "\033[33m", // 黄色
	ERROR: "\033[31m", // 红色
}

// resetColor 终端重置颜色
const resetColor = "\033[0m"

// InitLogger 初始化日志系统
func InitLogger(levelStr string, writer io.Writer) {
	levelStr = strings.ToLower(levelStr)
	level, ok := levelMap[levelStr]
	if !ok {
		level = INFO // 默认级别
	}

	currentLevel = level

	if writer == nil {
		writer = os.Stdout
	}

	logger = log.New(writer, "", log.LstdFlags|log.Lshortfile)
}

// internalLog 输出日志（根据级别过滤，并添加颜色）
func internalLog(level LogLevel, prefix string, msg string) {
	if level < currentLevel {
		return
	}

	color := colorMap[level]
	if color != "" {
		msg = fmt.Sprintf("%s%s%s", color, msg, resetColor)
	}

	logger.Output(3, fmt.Sprintf("[%s] %s", prefix, msg))
}

// 对外暴露的日志接口
func Debug(v ...any) {
	internalLog(DEBUG, "DEBUG", fmt.Sprint(v...))
}

func Info(v ...any) {
	internalLog(INFO, "INFO", fmt.Sprint(v...))
}

func Warn(v ...any) {
	internalLog(WARN, "WARN", fmt.Sprint(v...))
}

func Error(v ...any) {
	internalLog(ERROR, "ERROR", fmt.Sprint(v...))
}
