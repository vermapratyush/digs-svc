package logger

import (
	"github.com/astaxie/beego/logs"
)

var (
	commonLogger = logs.NewLogger(1)
)
func Initialize()  {
	//Configure logger
	commonLogger.EnableFuncCallDepth(true)
	commonLogger.SetLogger("multifile", `{"filename":"test.log","separate":["emergency", "alert", "critical", "error", "warning", "notice", "info", "debug"]}`)

}

func Debug(format string, v ...interface{}) {
	commonLogger.Debug(format, v)
}

func Error(format string, v ...interface{}) {
	commonLogger.Error(format, v)
}

func Critical(format string, v ...interface{}) {
	commonLogger.Critical(format, v)
}
