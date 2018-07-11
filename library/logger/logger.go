package logger

import (
	"log"
	"os"
	"io"
	"fmt"
	"whale-market/library/config"
)

var (
	Trace   *log.Logger // 记录所有日志
	Info    *log.Logger // 重要的信息
	Warning *log.Logger // 需要注意的信息
	Error   *log.Logger // 致命错误
)

func LoadCommonLogInit() {
	Trace = initCusLogger("trace.log", "TRACE: ")
	Info = initCusLogger("info.log", "INFO: ")
	Warning = initCusLogger("warning.log", "WARNING: ")
	Error = initCusLogger("error.log", "ERROR: ")
}

func initCusLogger(fileName, preFix string) (logger *log.Logger) {
	showInConsole := false
	confGet := config.GetConfigMap("env")
	logConf := confGet["log_console"].(string)
	if logConf == "true" {
		showInConsole = true
	}
	if confGet == nil {
		log.Fatalln("Failed to open error log file:", confGet)
		return nil
	}
	filePath := confGet["logpath"].(string)
	projectName := confGet["service_name"].(string)
	logpath := fmt.Sprintf("%s/%s.%s", filePath, projectName, fileName)
	if showInConsole {
		logger = log.New(io.Writer(os.Stdout), preFix, log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	} else {
		fileLog, err := os.OpenFile(logpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatalln("Failed to open error log file:", err)
			return nil
		}
		logger = log.New(io.Writer(fileLog), preFix, log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
	}
	return logger
}
