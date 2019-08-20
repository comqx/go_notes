package main

import (
	"time"

	"oldbody.com/day7/homework/loggergoroutine/mylogger"
)

// 1、支持往不同的地方输出日志
// 2、日志分级别
// Debug
// Trace
// Info
// Warning
// Error
// Fatal
// 3、日志要支持开关控制，比如说开发的时候什么级别都能输出，但是上线之后只有INFO级别往下的才能输出
// 4、完整的日志记录要包含有时间、行号、文件名、日志级别、日志信息
// 5、日志文件要切割
// 按文件大小切割
// 每次记录日志之前都判断一下当前写的这个文件的文件大小
// 按日期切割
// 在日志结构体中设置一个字段记录上一次切割的小时数
// 在写日志之前检查一下当前时间的小时数和之前保存的是否一致，不一致就要切割

func main() {
	// consoleLogger := mylogger.NewConsoleLogger("fatal")

	//级别、路径、文件、类型、大小、写文件线程数
	// fileLogger := mylogger.NewFileLogger("fatal", "./", "xx.log", "size", 5*1024, 5)
	fileLogger := mylogger.NewFileLogger("fatal", "./", "abc.log", "date", 1*10, 10)

	var logger mylogger.LoggerInterface
	logger = fileLogger
	// logger = consoleLogger

	name := "lqx"
	age := 20
	for {
		logger.Debug("name:%s,age:%d", name, age)
		logger.Info("name:%s,age:%d", name, age)
		logger.Warning("name:%s,age:%d", name, age)
		logger.Error("name:%s,age:%d", name, age)
		logger.Fatal("name:%s,age:%d", name, age)
		time.Sleep(time.Second * 1)
	}
}
