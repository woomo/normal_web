package logger

import (
	"fmt"
	"log"
	"normal_web/util"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	debugLogger   *log.Logger
	infoLogger    *log.Logger
	warnLogger    *log.Logger
	errorLogger   *log.Logger
	sqlLogger     *log.Logger
	logLevel      = 0
	logFile       string
	logout        *os.File
	day           int
	dayChangeLock sync.RWMutex
)

const (
	DebugLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

func Init(file string) {
	logConf := util.CreateConfig(file)
	SetLogLevel(logConf.GetString("level"))
	SetLogFile(logConf.GetString("file"))
}

func SetLogLevel(level string) {
	var l int
	switch level {
	case "DEBUG":
		l = DebugLevel
	case "INFO":
		l = InfoLevel
	case "WARM":
		l = WarnLevel
	case "ERROR":
		l = ErrorLevel
	}
	logLevel = l
}

func SetLogFile(file string) {
	logFile = file
	now := time.Now()
	var err error
	if logout, err = os.OpenFile(filepath.Join(util.ProjectRootPath, logFile), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o664); err == nil {
		debugLogger = log.New(logout, "[DEBUG]", log.LstdFlags)
		infoLogger = log.New(logout, "[INFO]", log.LstdFlags)
		warnLogger = log.New(logout, "[WARN]", log.LstdFlags)
		errorLogger = log.New(logout, "[ERROR]", log.LstdFlags)
		sqlLogger = log.New(logout, "[SOL]", log.LstdFlags)
		day = now.YearDay()
		dayChangeLock = sync.RWMutex{}
	} else {
		panic(err)
	}
}

func checkAndChangeLogFile() {
	dayChangeLock.Lock()
	defer dayChangeLock.Unlock()
	now := time.Now()
	if now.YearDay() == day {
		return
	}
	logout.Close()

	var err error
	postfix := now.Add(-24 * time.Hour).Format("20060102")
	if err = os.Rename(logFile, logFile+"-"+postfix); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("add date postpix %s to log file %s failed: %v\n", postfix, logFile, err))
		return
	}

	if logout, err = os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o664); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("open file %s failed: %v\n", logFile, err))
		return
	} else {
		debugLogger = log.New(logout, "[DEBUG]", log.LstdFlags)
		infoLogger = log.New(logout, "[INFO]", log.LstdFlags)
		warnLogger = log.New(logout, "[WARN]", log.LstdFlags)
		errorLogger = log.New(logout, "[ERROR]", log.LstdFlags)
		sqlLogger = log.New(logout, "[SOL]", log.LstdFlags)
		day = now.YearDay()
	}
}

func getFileAndLineOn() (string, string, int) {
	if funcName, file, line, ok := runtime.Caller(3); ok {
		return file, runtime.FuncForPC(funcName).Name(), line
	} else {
		return "", "", 0
	}
}

func addPrefix() string {
	file, _, line := getFileAndLineOn()
	arr := strings.Split(file, string(filepath.Separator))
	if len(arr) > 3 {
		arr = arr[len(arr)-3:]
	}
	return strings.Join(arr, string(filepath.Separator)) + ":" + strconv.Itoa(line)
}

func Debug(format string, v ...any) {
	if logLevel <= DebugLevel {
		checkAndChangeLogFile()
		debugLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Info(format string, v ...any) {
	if logLevel <= InfoLevel {
		checkAndChangeLogFile()
		infoLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Warn(format string, v ...any) {
	if logLevel <= WarnLevel {
		checkAndChangeLogFile()
		warnLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Error(format string, v ...any) {
	if logLevel <= ErrorLevel {
		checkAndChangeLogFile()
		errorLogger.Printf(addPrefix()+" "+format, v...)
	}
}

func Fatalf(format string, v ...any) {
	checkAndChangeLogFile()
	errorLogger.Fatalf(addPrefix()+" "+format, v...)
}
