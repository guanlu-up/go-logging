package logging

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"time"
)

type LogLevel uint16

type Logger struct {
	filepath string
	perm     os.FileMode
	Level    LogLevel
}

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

var LevelMapper = map[LogLevel]string{
	TRACE:   "TRACE",
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

func output(logger *Logger, level LogLevel, message string) {
	// 如果logger对象的输出等级大于传入的level则忽略这条日志
	if logger.Level > level {
		return
	}
	// Caller()的skip参数表示:需要向上冒泡的层级
	_, filepath, line, ok := runtime.Caller(2)
	if !ok {
		return
	}

	// 日志格式: [YY-mm-dd HH:MM:SS DEBUG] main.go:line: log message
	date := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(
		"[%s %s] %s:%d: %s\n",
		date, LevelMapper[level], path.Base(filepath), line, message)
	fmt.Print(msg)

	if logger.filepath != "" && logger.perm != 0x0 {
		mode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
		fio, err := os.OpenFile(logger.filepath, mode, logger.perm)
		if err != nil {
			// 如果打开文件出错则将相关属性值重制为空
			fmt.Println(err)
			logger.filepath = ""
			logger.perm = 0x0
			return
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(fio)

		_, err = fio.WriteString(msg)
	}
}

func (log *Logger) Debug(message string) {
	output(log, DEBUG, message)
}

func (log *Logger) Info(message string) {
	output(log, INFO, message)
}

func (log *Logger) Warning(message string) {
	output(log, WARNING, message)
}

func (log *Logger) Error(message string) {
	output(log, ERROR, message)
}

// NewLogger 创建一个新的Logger对象
// param level: 表示日志输出的等级, 低于此等级的日志则不被输出
// 日志等级包括: TRACE DEBUG INFO WARNING ERROR FATAL
func NewLogger(level LogLevel) *Logger {
	if level < TRACE || level > FATAL {
		message := fmt.Sprintf("param wrong, required: %d - %d", DEBUG, ERROR)
		panic(message)
	}
	return &Logger{Level: level}
}

// SetFile 设置记录日志的文件
// param filepath: 文件路径,需要为绝对路径
// param perm: 打开文件所需的权限,如 0644
func (log *Logger) SetFile(filepath string, perm os.FileMode) {
	log.filepath = filepath
	log.perm = perm
}
