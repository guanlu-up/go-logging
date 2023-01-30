package logging

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

type LogLevel uint16

type Logger struct {
	Level    LogLevel    // 日志输出的等级
	filepath string      // 存储日志的文件路径
	perm     os.FileMode // 文件自身的权限
	maxSize  uint64      // 文件最大可存储的容量
	fio      *os.File    // 日志文件打开后的对象
}

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

// LevelMapper  level to string
var LevelMapper = map[LogLevel]string{
	TRACE:   "TRACE",
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

// canBeWrite 判断当前日志文件是否可以被写入
// 当日志文件的容量大于Logger对象指定的最大容量值时就需要重新打开一个新文件
func canBeWrite(logger *Logger) bool {
	info, _ := logger.fio.Stat()
	if uint64(info.Size()) < logger.maxSize {
		return true
	}
	_ = logger.CloseFile()

	// 重命名原文件
	dir, filename := path.Split(logger.filepath)
	date := time.Now().Format("20060102150405")
	index := strings.LastIndex(filename, ".")
	backupFile := fmt.Sprintf("%s_%s%s", filename[:index], date, filename[index:])
	err := os.Rename(logger.filepath, path.Join(dir, backupFile))
	if err != nil {
		return false
	}

	// 重新创建原文件
	mode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	fio, err := os.OpenFile(logger.filepath, mode, logger.perm)
	if err != nil {
		return false
	}
	logger.fio = fio
	return true
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

	// 如果已经指定并打开了日志文件,则将日志一同记录到文件中
	if logger.fio != nil {
		if !canBeWrite(logger) {
			return
		}

		_, err := logger.fio.WriteString(msg)
		if err != nil {
			fmt.Println(err)
			err = logger.CloseFile()
		}
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
// 	param level: 表示日志输出的等级, 低于此等级的日志则不被输出
// 日志等级包括: TRACE DEBUG INFO WARNING ERROR FATAL
func NewLogger(level LogLevel) *Logger {
	if level < TRACE || level > FATAL {
		message := fmt.Sprintf("param wrong, required: %d - %d", DEBUG, ERROR)
		panic(message)
	}
	return &Logger{Level: level}
}

// SetFile 设置记录日志的文件
// 	param filepath: 文件路径,需要为完整路径
// 	param perm: 文件自身的权限,如 0644
// 	param maxSize: 文件最大可存储的容量,单位Bit
func (log *Logger) SetFile(filepath string, perm os.FileMode, maxSize uint64) {
	if maxSize <= 0 {
		panic("maxSize value invalid!, required > 0")
	}

	log.filepath = filepath
	log.perm = perm
	log.maxSize = maxSize

	mode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
	fio, err := os.OpenFile(log.filepath, mode, log.perm)
	if err != nil {
		panic(err)
	}
	log.fio = fio
}

func (log *Logger) CloseFile() error {
	if log.fio == nil {
		return fmt.Errorf("未指定日志文件,无法进行关闭")
	}
	err := log.fio.Close()
	log.fio = nil
	return err
}
