package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	"time"
)

//LogLevel 定义一个日志级别类型
type LogLevel uint8

//定义日志级别的常量
const (
	UNKNOWN LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

// LoggerInterface 定义了一个日志的接口，规范发送到文件和终端打印
type LoggerInterface interface {
	Debug(content string, logContent ...interface{})
	Info(content string, logContent ...interface{})
	Warning(content string, logContent ...interface{})
	Error(content string, logContent ...interface{})
	Fatal(content string, logContent ...interface{})
}

//打印出来执行caller的位置，（也就是打印出来打印日志的位置）
func runtimeCaller(n int) (string, string, int) {
	pc, file, line, ok := runtime.Caller(n)
	if !ok {
		fmt.Printf("runtime.Caller() failed\n")
		return "", "", 0
	}
	funcName := runtime.FuncForPC(pc).Name()
	filename := path.Base(file)
	return funcName, filename, line
}

//处理err问题 并且打印出来哪一行出现问题
func tryError(err error) {
	if err != nil {
		funcName, filename, line := runtimeCaller(2)
		fmt.Printf("[funcName:%s filename:%s line:%d] cont:%s", funcName, filename, line, err)
	}
}

//format 根据日志内容，做统一格式化
func contentFormat(content string, levelModeStr string, logContent ...interface{}) string {
	if logContent != nil {
		content = fmt.Sprintf(content, logContent...)
	}
	now := time.Now().Format("2006-01-02 15:04:05.000")
	funcName, filename, line := runtimeCaller(3)
	content = fmt.Sprintf("[%s] %s [funcName:%s filename:%s line:%d] %s", now, levelModeStr, funcName, filename, line, content)
	return content
}

// parseLogLevel 传入日志级别，然后解析成为系统定义的日志级别常量
func parseLogLevel(levelStr string) (LogLevel, error) {
	levelStr = strings.ToLower(levelStr)
	switch levelStr {
	case "debug", "Debug":
		return DEBUG, nil
	case "info", "Info":
		return INFO, nil
	case "warning", "Warning":
		return WARNING, nil
	case "error", "Error":
		return ERROR, nil
	case "fatal", "Fatal":
		return FATAL, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOWN, err
	}
}
