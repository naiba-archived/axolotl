package util

import (
	"fmt"
	"runtime"
	"time"
)

func Errorf(format string, v ...interface{}) {
	printLog("error", format, v...)
}

func Infof(format string, v ...interface{}) {
	printLog("info", format, v...)
}

func Debugf(format string, v ...interface{}) {
	printLog("debug", format, v...)
}

func printLog(level string, format string, v ...interface{}) {
	if pc, file, line, ok := runtime.Caller(2); ok {
		fName := runtime.FuncForPC(pc).Name()
		v = append([]interface{}{
			time.Now(), level, fName, file, line,
		}, v...)
		fmt.Printf("time=\"%s\" level=\"%s\" trace=\"%s %s:%d\" msg >> "+format, v...)
	}
}
