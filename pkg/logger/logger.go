// deprecated code
package logger

// import (
// 	"fmt"
// 	"os"
// 	"regexp"
// 	"strings"
// )

// var rep *regexp.Regexp = regexp.MustCompile(`\<.*\>`)
// var defaultKey = "default"
// var AppendToFile = false // if value True, log will append; False means override

// var DefaultLogger *Logger = nil

// type Logger struct {
// 	writeToStderr bool
// 	writeToFile   bool
// 	files         map[string]*Record
// }

// type Record struct {
// 	buffer []byte
// 	file   string
// }

// func NewDefaultLogger(writeToStderr, writeToFile bool, defaultPath string) *Logger {
// 	DefaultLogger = NewLogger(writeToStderr, writeToFile, defaultPath)
// 	return DefaultLogger
// }

// func NewLogger(writeToStderr, writeToFile bool, defaultPath string) *Logger {
// 	if defaultPath == "" {
// 		defaultPath = "unknown-log.log"
// 	}

// 	defaultRecord := &Record{
// 		buffer: make([]byte, 0, 1000),
// 		file:   defaultPath,
// 	}

// 	if !AppendToFile && writeToFile {
// 		if fp, err := os.OpenFile(defaultPath, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
// 			return nil
// 		} else {
// 			fp.Close()
// 		}
// 	}

// 	files := make(map[string]*Record)
// 	files[defaultKey] = defaultRecord
// 	return &Logger{
// 		writeToStderr: writeToStderr,
// 		writeToFile:   writeToFile,
// 		files:         files,
// 	}
// }

// func (logger *Logger) SetLog(key, logfile string) error {
// 	if logger == nil {
// 		return fmt.Errorf("nil logger")
// 	}

// 	if key != "" && logfile != "" {
// 		logger.files[key] = &Record{
// 			file:   logfile,
// 			buffer: make([]byte, 0, 1000),
// 		}
// 	}

// 	if !AppendToFile && logger.writeToFile {
// 		if fp, err := os.OpenFile(logfile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644); err != nil {
// 			return nil
// 		} else {
// 			fp.Close()
// 		}
// 	}

// 	return fmt.Errorf("bad args")
// }

// func (logger *Logger) ClearBuffer() {
// 	if logger == nil {
// 		return
// 	}

// 	for _, record := range logger.files {
// 		record.buffer = record.buffer[0:0] // only clear
// 	}
// }

// func (logger *Logger) RemoveLog(key string) {
// 	if logger == nil {
// 		return
// 	}
// 	delete(logger.files, key)
// }

// // Implement io.Writer interface
// func (logger *Logger) Write(p []byte) (n int, err error) {
// 	if logger == nil {
// 		return 0, fmt.Errorf("nil logger")
// 	}
// 	if p == nil {
// 		return 0, nil
// 	}

// 	lines := strings.Split(string(p), "\n")

// 	length := 0
// 	for _, line := range lines {
// 		if len(line) < 1 {
// 			continue
// 		}

// 		key, data := findAndReplace(line)

// 		length += len(data)
// 		if logger.writeToStderr {
// 			fmt.Println(string(data))
// 		}

// 		if logger.writeToFile {
// 			if _, ok := logger.files[key]; ok {
// 				logger.files[key].buffer = append(append(logger.files[key].buffer, data...), []byte("\n")...)
// 			} else {
// 				logger.files[defaultKey].buffer = append(append(logger.files[defaultKey].buffer, data...), []byte("\n")...)
// 			}
// 		}
// 	}

// 	if logger.writeToFile {
// 		if err := logger.WriteBuffer(); err != nil {
// 			return 0, err
// 		}
// 	}

// 	return length, nil
// }

// func (logger *Logger) WriteBuffer() error {
// 	if logger == nil {
// 		return fmt.Errorf("nil logger")
// 	}

// 	defer logger.ClearBuffer()
// 	for _, record := range logger.files {
// 		f, err := os.OpenFile(record.file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 		if err != nil {
// 			return err
// 		}

// 		if _, err := f.Write(record.buffer); err != nil {
// 			return err
// 		}

// 		if err := f.Close(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// // find rep-string and delete in raw-string
// // return rep-string-$1, new-string
// func findAndReplace(raw string) (string, []byte) {
// 	datas := []byte(raw)
// 	indexs := rep.FindIndex([]byte(raw))
// 	if indexs == nil {
// 		return "", []byte(raw)
// 	}

// 	s := string(datas[indexs[0]+1 : indexs[1]-1])

// 	//return s, datas
// 	return s, append(datas[:indexs[0]], datas[indexs[1]:]...)
// }
