package util

import (
	"fmt"
	"runtime"
	"time"
)

func Errorf(callStackDepth int, format string, v ...interface{}) {
	printLog("error", 1+callStackDepth, format, v...)
}

func Infof(callStackDepth int, format string, v ...interface{}) {
	printLog("info", 1+callStackDepth, format, v...)
}

func Debugf(callStackDepth int, format string, v ...interface{}) {
	printLog("debug", 1+callStackDepth, format, v...)
}

func printLog(level string, callStackDepth int, format string, v ...interface{}) {
	if pc, file, line, ok := runtime.Caller(1 + callStackDepth); ok {
		fName := runtime.FuncForPC(pc).Name()
		v = append([]interface{}{
			time.Now(), level, fName, file, line,
		}, v...)
		fmt.Printf("time=\"%s\" level=\"%s\" trace=\"%s %s:%d\" "+format, v...)
	}
}
