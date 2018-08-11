package logger

import (
	"log"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"encoding/json"
	"account/library/config"
)

var (
	ZapTrace   *zap.Logger // 记录所有日志
	ZapInfo    *zap.Logger // 重要的信息
	ZapWarning *zap.Logger // 需要注意的信息
	ZapError   *zap.Logger // 致命错误
)

func initLogger(fileName string) *zap.Logger {
	var configTemp string
	showInConsole := false
	confGet := config.GetConfigMap("env")
	filePath := confGet["logpath"].(string)
	if confGet != nil {
		logConf := confGet["log_console"].(string)
		if logConf == "true" {
			showInConsole = true
		}
	}
	projectName := confGet["service_name"].(string)
	path := fmt.Sprintf("%s/%s.%s", filePath, projectName, fileName)
	if showInConsole {
		configTemp = fmt.Sprintf(`{
      "level": "DEBUG",
      "encoding": "json",
      "outputPaths": ["stdout"],
      "errorOutputPaths": ["stdout"]
      }`)
	} else {
		configTemp = fmt.Sprintf(`{
      "level": "DEBUG",
      "encoding": "json",
      "outputPaths": ["%s"],
      "errorOutputPaths": ["%s"]
      }`, path, path)
	}
	var cfg zap.Config
	if err := json.Unmarshal([]byte(configTemp), &cfg); err != nil {
		panic(err)
	}
	cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := cfg.Build()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile | log.LstdFlags)
	return logger
}

func init() {
	ZapInfo = initLogger("info.log")
	ZapTrace = initLogger("trace.log")
	ZapWarning = initLogger("warning.log")
	ZapError = initLogger("error.log")
}