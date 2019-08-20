package mylogger

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"
)

// FileLogger 定义了一个filelogger的结构体
type FileLogger struct {
	Level           LogLevel
	filePath        string
	fileName        string
	fileSegType     string
	fileSize        int64
	fileDateSeg     int64
	fileDateLast    time.Time
	fileDateLastErr time.Time
	fileObj         *os.File
	fileObjErr      *os.File
}

// loggerTxt 定义日志内容的结构体
type loggerTxt struct {
	loggerContent string
}

// 定义俩个通道，用来接收日志
var chcontent = make(chan *loggerTxt, 256)
var chcontentErr = make(chan *loggerTxt, 256)

//NewFileLogger 制作了一个filelogger的构造函数,开启俩个线程读取chan里面的日志，然后往文件中写入日志
func NewFileLogger(levelStr, filePath, fileName, fileSegTypePri string, fileNum int64, threadNum int) *FileLogger {
	logLevel, err := parseLogLevel(levelStr)
	tryError(err)
	fileSegType := strings.ToLower(fileSegTypePri)

	//实例化一个对象
	fl := &FileLogger{
		Level:       logLevel,
		filePath:    filePath,
		fileName:    fileName,
		fileSegType: fileSegType,
	}
	if fileSegType == "size" {
		fl.fileSize = fileNum
	} else if fileSegType == "date" {
		fl.fileDateSeg = fileNum

	} else {
		err := errors.New("无效的大小，如果是按照大小切分日志，那么请输入1*1024*1024格式，如果是按照时间切分日志，那么请输入1*60*60*60格式。")
		tryError(err)
		os.Exit(1)
	}

	//初始化文件具柄
	fl.initFile()

	//开启10个线程写普通日志的线程和5个写err日志的线程
	for i := 1; i <= threadNum; i++ {
		//开启后台往文件中写入的goroutine，监听通道里面的数据
		go func(chcontent <-chan *loggerTxt) {
			for content := range chcontent {
				fmt.Fprintln(fl.fileObj, content.loggerContent)
			}
		}(chcontent)

		// 开启后台往文件中写入的goroutine，监听通道里面的数据，专门用来写err日志
		go func(chcontentErr <-chan *loggerTxt) {
			for content := range chcontentErr {
				fmt.Fprintln(fl.fileObjErr, content.loggerContent)
			}
		}(chcontentErr)
	}

	return fl
}

//initFile 创建文件，并初始化文件具柄(这里必须使用指针模式，因为涉及到对FileLogger的属性进行赋值了)
func (f *FileLogger) initFile() {
	filePathName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(filePathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	now := time.Now()
	tryError(err)
	if f.Level >= ERROR {
		filePathName = path.Join(f.filePath, "error.log")
		fileObjErr, err := os.OpenFile(filePathName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		tryError(err)
		f.fileObjErr = fileObjErr
		f.fileDateLastErr = now
	}
	f.fileObj = fileObj
	f.fileDateLast = now
}

//enable 根据日志级别，启动对应的级别日志
func (f *FileLogger) enable(levelmode LogLevel) bool {
	//f.level是前端传过来需要打印的日志级别 INFO 2
	// levelmode是各个级别传过来的自己的级别 debug 1
	if f.Level >= levelmode {
		return true
	}
	return false
}

//检测文件大小，或者检测文件创建的时间 然后执行重命名文件函数
func (f *FileLogger) checkRenameFile() {
	now := time.Now()
	fileStat, err := f.fileObj.Stat()
	tryError(err)
	fileStatErr, err := f.fileObjErr.Stat()
	tryError(err)
	dateLastSeg := int64(now.Sub(f.fileDateLast).Seconds())
	// dateLastSegErr := int64(now.Sub(f.fileDateLastErr).Seconds())
	if f.fileSegType == "size" {
		if fileStat.Size() >= f.fileSize {
			f.fileObj.Close()
			f.RenameFileOpenFile(1, 0)
		} else if fileStatErr.Size() >= f.fileSize {
			f.fileObjErr.Close()
			f.RenameFileOpenFile(2, 0)
		}
	}
	fmt.Println("fileSegType:", f.fileSegType)
	if f.fileSegType == "date" {
		fmt.Printf("已经跑了的时间 ：%d  时间间隔：%d\n", dateLastSeg, f.fileDateSeg)
		if dateLastSeg >= f.fileDateSeg {
			f.fileObj.Close()
			f.fileObjErr.Close()
			f.RenameFileOpenFile(1, 2)
		}
	}
}

// RenameFileOpenFile 重命名文件，然后新建一个文件
func (f *FileLogger) RenameFileOpenFile(flagNum1, flagNum2 int) {
	bakTime := time.Now().Format("20060102150405")
	if flagNum2 == 0 {
		flagNum2 = flagNum1
	}
	if flagNum1 == 1 {
		filePathName := path.Join(f.filePath, f.fileName)
		newFilePathName := fmt.Sprintf("%s-%s", filePathName, bakTime)
		err := os.Rename(filePathName, newFilePathName)
		fmt.Println(filePathName, newFilePathName)
		tryError(err)
		f.initFile()
	}
	if flagNum2 == 2 {
		filePathNameErr := path.Join(f.filePath, "error.log")
		newFilePathNameErr := fmt.Sprintf("%s-%s", filePathNameErr, bakTime)
		err := os.Rename(filePathNameErr, newFilePathNameErr)
		fmt.Println(filePathNameErr, newFilePathNameErr)
		tryError(err)
		f.initFile()
	}
}

// Debug loglevel debug
func (f *FileLogger) Debug(content string, logContent ...interface{}) {
	if f.enable(DEBUG) {
		f.checkRenameFile()
		logContentTxt := loggerTxt{
			loggerContent: contentFormat(content, "DEBUG", logContent...),
		}
		//使用select,把日志的指针写入chan,如果管道满了，就主动丢弃日志
		select {
		case chcontent <- &logContentTxt:
		default:
			//丢弃当前日志
			time.Sleep(time.Millisecond * 300)
		}
	}

}

// Info loglevel info
func (f *FileLogger) Info(content string, logContent ...interface{}) {
	if f.enable(INFO) {
		f.checkRenameFile()
		logContentTxt := loggerTxt{
			loggerContent: contentFormat(content, "INFO", logContent...),
		}
		//使用select,把日志的指针写入chan,如果管道满了，就主动丢弃日志
		select {
		case chcontent <- &logContentTxt:
		default:
			//丢弃当前日志
			time.Sleep(time.Millisecond * 300)
		}

	}
}

// Warning loglevel warning
func (f *FileLogger) Warning(content string, logContent ...interface{}) {
	if f.enable(WARNING) {
		f.checkRenameFile()
		logContentTxt := loggerTxt{
			loggerContent: contentFormat(content, "WARNING", logContent...),
		}
		//使用select,把日志的指针写入chan,如果管道满了，就主动丢弃日志
		select {
		case chcontent <- &logContentTxt:
		default:
			//丢弃当前日志
			time.Sleep(time.Millisecond * 300)
		}
	}
}

// Error loglevel error
func (f *FileLogger) Error(content string, logContent ...interface{}) {
	if f.enable(ERROR) {
		f.checkRenameFile()
		logContentTxt := loggerTxt{
			loggerContent: contentFormat(content, "ERROR", logContent...),
		}
		//使用select,把日志的指针写入chan,如果管道满了，就主动丢弃日志
		select {
		case chcontentErr <- &logContentTxt:
		default:
			//丢弃当前日志
			time.Sleep(time.Millisecond * 300)
		}
	}
}

// Fatal loglevel fatal
func (f *FileLogger) Fatal(content string, logContent ...interface{}) {
	if f.enable(FATAL) {
		f.checkRenameFile()
		logContentTxt := loggerTxt{
			loggerContent: contentFormat(content, "FATAL", logContent...),
		}
		//使用select,把日志的指针写入chan,如果管道满了，就主动丢弃日志
		select {
		case chcontentErr <- &logContentTxt:
		default:
			//丢弃当前日志
			time.Sleep(time.Millisecond * 300)
		}
	}
}
