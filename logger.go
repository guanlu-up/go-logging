package logging

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	filepath string
	perm     os.FileMode
}

const (
	DEBUG   = "debug"
	INFO    = "info"
	WARNING = "warning"
	ERROR   = "error"
)

func output(logger Logger, level string, message string) {
	date := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf("\n[%v %s]:%s", date, level, message)
	_, _ = os.Stdin.WriteString(msg)

	if logger.filepath != "" && logger.perm != 0x0 {
		mode := os.O_WRONLY | os.O_CREATE | os.O_APPEND
		fileIO, err := os.OpenFile(logger.filepath, mode, logger.perm)
		if err != nil {
			return
		}
		defer func(io *os.File) {
			_ = io.Close()
		}(fileIO)

		_, writeError := fileIO.WriteString(msg)
		if writeError != nil {
			return
		}
	}
}

func (log Logger) Debug(message string) {
	output(log, DEBUG, message)
}

func (log Logger) Info(message string) {
	output(log, INFO, message)
}

func (log Logger) Warning(message string) {
	output(log, WARNING, message)
}

func (log Logger) Error(message string) {
	output(log, ERROR, message)
}

func NewLogger(filepath string, perm os.FileMode) *Logger {
	return &Logger{filepath: filepath, perm: perm}
}
