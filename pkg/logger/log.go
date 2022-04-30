package logger

import (
	"io"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

type LogWriter struct {
	WriteToStderr bool
	WriteToFile   bool
	LogPath       string
	LogWriter     io.Writer
}

func NewLogWriter(writeToStderr, writeToFile bool, defaultPath string) *LogWriter {
	if defaultPath == "" {
		defaultPath = "unknown-log.log"
	}

	defaultPath = strings.SplitN(defaultPath, ".log", 2)[0]

	lw := &LogWriter{
		WriteToStderr: writeToStderr,
		WriteToFile:   writeToFile,
		LogPath:       defaultPath,
		LogWriter:     nil,
	}

	writerList := make([]io.Writer, 0, 2)
	if lw.WriteToStderr {
		writerList = append(writerList, os.Stdout)
	}

	if lw.WriteToFile {
		if fw, err := rotatelogs.New(
			lw.LogPath+"_%Y-%m-%d.log",
			rotatelogs.WithMaxAge(time.Duration(168)*time.Hour), // 7 days circle
			rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
		); err != nil {
			panic(err)
		} else {
			writerList = append(writerList, fw)
		}

	}

	if len(writerList) > 0 {
		lw.LogWriter = io.MultiWriter(writerList...)
	}

	return lw
}

func (lw *LogWriter) Logger() io.Writer {
	if lw == nil {
		return nil
	}
	return lw.LogWriter
}
