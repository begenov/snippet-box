package main

import (
	"log"
	"os"
)

// Loggers
var (
	logInfo *log.Logger
	logErr  *log.Logger
	logWarn *log.Logger
)

const filename = "../log_file.log"

type IAuthLoggers interface {
	Info(v ...interface{})
	Err(v ...interface{})
	Warn(v ...interface{})
}

type Loggers struct {
	logInfo *log.Logger
	logWarn *log.Logger
	logErr  *log.Logger
}

func init() {
	flags := log.LstdFlags | log.Lshortfile
	fileInfo, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileWarn, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	fileErr, _ := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	logInfo = log.New(fileInfo, "INFO:\t", log.LstdFlags)
	logWarn = log.New(fileWarn, "WARN:\t", flags)
	logErr = log.New(fileErr, "ERR:\t", flags)
	// Set the logger to write to the multiwriter.

}

func NewLogger() IAuthLoggers {
	return &Loggers{
		logInfo: logInfo,
		logWarn: logWarn,
		logErr:  logErr,
	}
}

func (l *Loggers) Info(v ...interface{}) {
	l.logInfo.Println(v...)
}

func (l *Loggers) Err(v ...interface{}) {
	l.logErr.Println(v...)
}

func (l *Loggers) Warn(v ...interface{}) {
	l.logWarn.Println(v...)
}
