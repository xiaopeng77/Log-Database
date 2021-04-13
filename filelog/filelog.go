package filelog

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

//定义6个级别的常量1.debug 2.Trace 3.Info 4.Warning  5.Error  6.Fatal
const (
	//定义日志级别
	UNKOWN LogLevel = iota
	Debug
	Trace
	Info
	Warning
	Error
	Fatal
)

type LogLevel uint16

//定义一个向文件写日志的结构体
type Logger struct {
	Level       LogLevel
	filePath    string //文件路径
	fileName    string //文件名
	fileMaxsize int64
}

//判断日志级别
func parseLoglevel(levelStr string) (LogLevel, error) {
	levelStr = strings.ToLower(levelStr)
	switch levelStr {
	case "debug":
		return Debug, nil
	case "trace":
		return Trace, nil
	case "info":
		return Info, nil
	case "warning":
		return Warning, nil
	case "error":
		return Error, nil
	case "fatal":
		return Fatal, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKOWN, err
	}
}

//将LogLevel类型的日志级别转换为String类型
func getLogString(l LogLevel) string {
	switch l {
	case 1:
		return "Debug"
	case 2:
		return "Trace"
	case 3:
		return "Info"
	case 4:
		return "Warning"
	case 5:
		return "Error"
	case 6:
		return "Fatal"
	default:
		return "UNKOWN"
	}
}

//Newlog 构造函数
func Newlog(levelStr string, filePath string, fileName string, fileMaxsize int64) Logger { //, filePath string, fileName string
	level, err := parseLoglevel(levelStr)
	if err != nil {
		panic(err)
	}
	return Logger{
		Level:       level,
		filePath:    filePath, //文件路径
		fileName:    fileName, //文件名
		fileMaxsize: fileMaxsize,
	}
}

//输出日志信息当文件大小超出设定大小时则新建一个日志文件来存放日志
func (l Logger) pringfLog(msg string, log LogLevel) {
	if log >= l.Level {
		//利用runtime.Caller()方法获取执行代码时的，函数名，包名，行号
		pc, file, line, ok := runtime.Caller(2) //参数为0时获取当前位置，参数为1时向上获取1层调用该函数的位置，以此类推
		if !ok {
			fmt.Println("runtiome.Caller() failed")
			return
		}
		funcName := runtime.FuncForPC(pc).Name() //获取函数名
		now := time.Now()
		//打开文件
		fileObj, err := os.OpenFile(l.filePath+l.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 999)
		if err != nil {
			fmt.Printf("Open file failed,err:%v\n", err)
			return
		}
		fileInfo, err := fileObj.Stat()
		if err != nil {
			fmt.Printf("Gain fileSize failed,err:%v\n", err)
			return
		}
		//fmt.Println(fileInfo.Size()) //查看一条日志信息为多少
		//判断日志文件剩下的空间能否写下一条日志信息如果不能则新建一个日志文件用来存放日志信息
		if (fileInfo.Size() + int64(150)) > l.fileMaxsize {
			fileObj.Close()
			//重名名当前日志文件
			os.Rename(l.filePath+l.fileName, l.filePath+l.fileName+now.Format("20060102150405"))
			//打开一个新的日志文件
			fileObj, err = os.OpenFile(l.filePath+l.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 999)
			if err != nil {
				fmt.Printf("Open file failed,err:%v\n", err)
				return
			}
		}
		fmt.Fprintf(fileObj, "[%s][%s][%s:%s:%d][%s]\n", now.Format("2006/01/02 15:04:05"), getLogString(log), file, funcName, line, msg)
		if log >= Error {
			//在当前目录下打开一个存放error和Fatal等级的日志文件
			errfileObj, err := os.OpenFile(l.filePath+l.fileName+"err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 999)
			if err != nil {
				fmt.Printf("Open file failed,err:%v\n", err)
				return
			}
			errfileInfo, err := errfileObj.Stat()
			if err != nil {
				fmt.Printf("Gain fileSize failed,err:%v\n", err)
				return
			}
			//fmt.Println(errfileInfo.Size()) //查看一条日志信息为多少
			//判断日志文件剩下的空间能否写下一条日志信息如果不能则新建一个日志文件用来存放日志信息
			if (errfileInfo.Size() + int64(150)) > l.fileMaxsize {
				errfileObj.Close()
				//重命名当前日志文件
				os.Rename(l.filePath+l.fileName+"err", l.filePath+l.fileName+now.Format("20060102150405")+"err")
				//打开一个新的日志文件
				errfileObj, err = os.OpenFile(l.filePath+l.fileName+"err", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 999)
				if err != nil {
					fmt.Printf("Open file failed,err:%v\n", err)
					return
				}
			}
			fmt.Fprintf(errfileObj, "[%s][%s][%s:%s:%d][%s]\n", now.Format("2006/01/02 15:04:05"), getLogString(log), file, funcName, line, msg)
			errfileObj.Close()
		}
		fileObj.Close()

	}
}

//Debug...
func (l Logger) Debug(msg string) {
	l.pringfLog(msg, Debug)
}
func (l Logger) Trace(msg string) {
	l.pringfLog(msg, Trace)
}
func (l Logger) Info(msg string) {
	l.pringfLog(msg, Info)
}
func (l Logger) Warning(msg string) {
	l.pringfLog(msg, Warning)
}
func (l Logger) Error(msg string) {
	l.pringfLog(msg, Error)
}
func (l Logger) Fatal(msg string) {
	l.pringfLog(msg, Fatal)
}
