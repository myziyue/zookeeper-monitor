package utils

import (
	"log"
	"os"
	"time"
)

const (
	TRACE     int = 1
	DEBUG     int = 100
	INFO      int = 200
	NOTICE    int = 250
	WARNING   int = 300
	ERROR     int = 400
	CEITICAL  int = 500
	ALERT     int = 550
	EMERGENCY int = 600

	Level_TRACE     string = "TRACE"
	Level_DEBUG     string = "DEBUG"
	Level_INFO      string = "INFO"
	Level_NOTICE    string = "NOTICE"
	Level_WARNING   string = "WARNING"
	Level_ERROR     string = "ERROR"
	Level_CEITICAL         = "CEITICAL"
	Level_ALERT            = "ALERT"
	Level_EMERGENCY        = "EMERGENCY"
)

var (
	logLevel, _ = GetOption("log_level", "system")
	loger       *log.Logger
)

func GetLevel(levelStr string) int {
	level := INFO
	switch levelStr {
	case Level_DEBUG:
		level = DEBUG
		break
	case Level_INFO:
		level = INFO
		break
	case Level_NOTICE:
		level = NOTICE
		break
	case Level_WARNING:
		level = WARNING
		break
	case Level_ERROR:
		level = ERROR
		break
	case Level_CEITICAL:
		level = CEITICAL
		break
	case Level_ALERT:
		level = ALERT
		break
	case Level_EMERGENCY:
		level = EMERGENCY
		break
	case Level_TRACE:
		level = TRACE
		break
	}
	return level
}

func init() {
	file := "logs/" + time.Now().Format("2006-01-02") + ".txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err)
	}
	loger = log.New(logFile, "[Zookeeper Monitor]", log.LstdFlags|log.Lshortfile|log.LUTC) // 将文件设置为loger作为输出
	return
}

func logger(content string, level string) {
	if logLevel == "" {
		Warning("config item 'log_level' not setting.")
		logLevel = Level_INFO
	}

	if GetLevel(level) < GetLevel(logLevel) {
		return
	}

	loger.Println("[", level, "] ", content) // log 还是可以作为输出的前缀
	log.Println("[", level, "] ", content)
}

func Debug(content string) {
	logger(content, Level_DEBUG)
}

func Info(content string) {
	logger(content, Level_INFO)
}

func Notice(content string) {
	logger(content, Level_NOTICE)
}

func Warning(content string) {
	logger(content, Level_WARNING)
}

func Errors(content string) {
	logger(content, Level_ERROR)
}

func Ceitical(content string) {
	logger(content, Level_CEITICAL)
}

func Alert(content string) {
	logger(content, Level_ALERT)
}

func Emergency(content string) {
	logger(content, Level_EMERGENCY)
}

func Trace(content string) {
	logger(content, Level_TRACE)
}
