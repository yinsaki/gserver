/*
@Time : 2018/6/4 20:18
@Author : yinsaki
@File : fileutil
*/

package log

import (
	"time"
	"sync"
	"os"
	"log"
	"yinsaki/gserver/core/system"
	"strconv"
	"fmt"
	"strings"
	"runtime"
)

type UNIT int64

const (
	_ = iota
	KB UNIT = 1 << (iota *10)
	MB
	GB
	TB
)

type LogLevel int
/*
@Time : 2018/6/5 10:21
@Author : yinsaki
@File : fileutil
*/

const (
	INFO LogLevel = iota
	DEBUG
	WARN
	ERROR
	FATAL
)

func (this LogLevel) String() string {
	switch this {
	case INFO :
		return "[INFO]"
	case DEBUG :
		return "[DBUG]"
	case WARN :
		return "[WARN]"
	case ERROR :
		return "[ERRO]"
	case FATAL :
		return "[FATL]"
	default :
		return "[DBUG]"
	}
}

type logFile struct {
	dir string
	filename string
	_suffix int
	isCover bool
	_date *time.Time
	mu *sync.RWMutex
	logFile *os.File
	lg *log.Logger
}

var (
	// log file
	logLevel LogLevel = INFO
	maxFileSize int64
	maxFileCount int32
	dailyCount int32
	consoleAppender bool = true
	dailyRolling bool = true
	fileRolling bool = false
	logObj *logFile

	consoleFormat string = "%s:%d %s %s"
	logFormat string = "%s %s"
)

const DATEFORMAT = "2018-06-05"

func (this *logFile) isMustRename() bool {
	if dailyRolling {
		t, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
		if t.After(*this._date) {
			return true
		}
	}

	if maxFileCount > 1 && system.GetFileSize(this.dir + "/" + this.filename) >= maxFileSize {
		return true
	}

	return false
}

func (this *logFile) rename() {
	if dailyRolling {
		this.dailyRollRename()
	} else {
		this.fileRollRename()
	}
}
func (this *logFile)dailyRollRename() {
	if this.logFile != nil {
		this.logFile.Close()
	}

	mainFileName := this.dir + "/" + this.filename
	fileName := this.dir + "/" + this.filename + this._date.Format(DATEFORMAT)
	err := os.Rename(mainFileName, fileName)
	if err != nil {
		this.lg.Println("file rename err", fileName, err.Error())
	}
	tt, _ := time.Parse(DATEFORMAT, time.Now().Format(DATEFORMAT))
	this._date = &tt
	this.logFile, _ = os.Create(mainFileName)
	this.lg = log.New(logObj.logFile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
}

func (this *logFile) nextSuffix() int {
	return int(this._suffix % int(maxFileSize) + 1)
}

func (this *logFile)fileRollRename() {
	this._suffix = this.nextSuffix()
	if this.logFile != nil {
		this.logFile.Close()
	}

	mainFileName := this.dir + "/" + this.filename
	fileName := this.dir + "/" + this.filename + strconv.Itoa(this._suffix)
	if system.IsFileExit(fileName) {
		os.Remove(fileName)
	}
	os.Rename(mainFileName, fileName)
	this.logFile, _ = os.Create(mainFileName)
	this.lg = log.New(logObj.logFile, "\n", log.Ldate|log.Ltime|log.Lshortfile)
}

func fileMonitor() {
	timer := time.NewTicker(1 * time.Second)
	for {
		select {
			case <- timer.C:
				fileCheck()
		}
	}
}

func fileCheck() {
	defer system.CatchError()

	if logObj != nil && logObj.isMustRename() {
		logObj.mu.Lock()
		defer  logObj.mu.Unlock()
		logObj.rename()
	}
}

// --------------------runtime-----------------------------
var curStackFlag bool = false
var curStackPath string
var curStackLine int

func getCurStackPath() string {
	if !curStackFlag {
		curStackFlag = true
		_, curStackPath, curStackLine, _ = runtime.Caller(0)
	}
	return curStackPath
}

func detectStack() (string, int) {
	curPath := getCurStackPath()
	for skip:= 0; ; skip ++ {
		_, path, line, ok := runtime.Caller(skip)
		if path != curPath {
			return path, line
		}
		if !ok {
			break
		}
	}
	return "", 0
}

func getTraceDirInfo(dir string) string {
	if system.GetOsFlag() == system.OS_WIN {
		split := strings.Split(dir, "\\")
		if len(split) > 2 {
			return split[0] + "\\" + split[1] + "\\...\\" + split[len(split)-1] + "\\"
		} else {
			return dir + "\\"
		}
	}

	split := strings.Split(dir, "/")
	if len(split) > 2 {
		return split[0] + "/.../" + split[len(split) - 1]
	} else {
		return dir + "/"
	}
}

func getTraceFileLine() (string, int) {
	path, line := detectStack()
	dir, file := system.SplitDirFile(path)
	dir = getTraceDirInfo(dir)
	return dir+file, line
}
//------------------------------------------------------------

func console(msg string) {
	if consoleAppender {
		log.Print(msg)
	}
}

func buildLogMessage(level LogLevel, msg string) string{
	file, line := getTraceFileLine()
	return fmt.Sprintf(consoleFormat + system.GetOsEol(), file, line, level.String(), msg)
}

func Trace(level LogLevel, format string, v ... interface{})bool {
	if dailyRolling {
		fileCheck()
	}

	defer system.CatchError()
	logObj.mu.Lock()
	defer logObj.mu.Unlock()

	msg := fmt.Sprintf(format, v)
	logMsg := buildLogMessage(level, msg)
	console(logMsg)

	if level >= logLevel {
		logObj.lg.Output(2, logMsg)
	}

	return true
}

func Info(format string, v ...interface{}) bool {
	return Trace(INFO, format, v)
}

func Debug(format string, v ...interface{}) bool {
	return Trace(DEBUG, format, v)
}

func Warn(format string, v ...interface{}) bool {
	return Trace(WARN, format, v)
}

func Error(format string, v ...interface{}) bool {
	return Trace(ERROR, format, v)
}

func Fatal(format string, v ...interface{}) bool {
	return Trace(FATAL, format, v)
}