package mylogger

import (
	"fmt"
)

// ConsoleLogger 定义了一个日志级别的struct
type ConsoleLogger struct {
	Level LogLevel
}

//NewConsoleLogger 定义一个构造函数
func NewConsoleLogger(levelStr string) ConsoleLogger {
	level, err := parseLogLevel(levelStr)
	tryError(err)
	return ConsoleLogger{
		Level: level,
	}
}

//enable 根据日志级别，启动对应的级别日志
func (c ConsoleLogger) enable(levelmode LogLevel) bool {
	//c.level是前端传过来需要打印的日志级别 INFO 2
	// levelmode是各个级别传过来的自己的级别 debug 1
	if c.Level >= levelmode {
		return true
	}
	return false
}

// Debug loglevel debug
func (c ConsoleLogger) Debug(content string, logContent ...interface{}) {
	if c.enable(DEBUG) {
		fmt.Println(contentFormat(content, "DEBUG", logContent...))
	}

}

// Info loglevel info
func (c ConsoleLogger) Info(content string, logContent ...interface{}) {
	if c.enable(INFO) {
		fmt.Println(contentFormat(content, "INFO", logContent...))
	}
}

// Warning loglevel warning
func (c ConsoleLogger) Warning(content string, logContent ...interface{}) {
	if c.enable(WARNING) {
		fmt.Println(contentFormat(content, "WARNING", logContent...))
	}
}

// Error loglevel error
func (c ConsoleLogger) Error(content string, logContent ...interface{}) {
	if c.enable(ERROR) {
		fmt.Println(contentFormat(content, "ERROR", logContent...))
	}
}

// Fatal loglevel fatal
func (c ConsoleLogger) Fatal(content string, logContent ...interface{}) {
	if c.enable(FATAL) {
		fmt.Println(contentFormat(content, "FATAL", logContent...))
	}
}
