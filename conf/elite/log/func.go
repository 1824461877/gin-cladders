package log

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
)

type FileFunc struct{
	level int
	logLevelFloor int
	logLevelCeil int
	file *os.File
	fileBufWriter *bufio.Writer
}

type FileWriteContainer struct {
	Write Write
	Bool bool
}

func CreateLogFile(filename string) (*FileFunc,error) {
	ff := &FileFunc{}
	if err := os.MkdirAll(path.Dir(filename), 0755); err != nil {
		if !os.IsExist(err) {
			return nil,err
		}
	}

	if file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return nil,err
	} else {
		ff.file = file
	}

	// w.fileBufWrite
	if ff.fileBufWriter = bufio.NewWriterSize(ff.file, 8192); ff.fileBufWriter == nil {
		return nil,errors.New("new fileBufWriter failed.")
	}

	return ff,nil
}


func WriteFileBufWrite(w *FileFunc,r *Record) error{
	level := r.level
	if level < w.logLevelFloor || level > w.logLevelCeil {
		return nil
	}
	if _, err := w.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}

func Flush(w *FileFunc) error {
	if w.fileBufWriter != nil {
		return w.fileBufWriter.Flush()
	}
	return nil
}

func (w *FileFunc) SetLogLevelFloor(floor int) {
	w.logLevelFloor = floor
}

func (w *FileFunc) SetLogLevelCeil(ceil int) {
	w.logLevelCeil = ceil
}

type ConsoleFunc struct {
	color bool
}

func (w *ConsoleFunc) Write(caw *ConsoleAndWrite,r *Record) error {
	if w.color {
		fmt.Fprint(os.Stdout, ((*colorRecord)(r)).OutputTemplate())
	} else {
		fmt.Fprint(os.Stdout, r.String())
	}
	return nil
}

var (
	logger_default *Logger
	takeup         = false
	config *LogConfig
)


func Trace(fmt string, args ...interface{}) {
	defaultLoggerInit(config)
	logger_default.deliverRecordToWriter(TRACE, fmt, args...)
}

func Debug(fmt string, args ...interface{}) {
	defaultLoggerInit(config)
	logger_default.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	defaultLoggerInit(config)
	logger_default.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	defaultLoggerInit(config)
	logger_default.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	defaultLoggerInit(config)
	logger_default.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	defaultLoggerInit(config)
	logger_default.deliverRecordToWriter(FATAL, fmt, args...)
}

//func Register(w Writer) {
//	defaultLoggerInit()
//	logger_default.Register(w)
//}

func Close() {
	defaultLoggerInit(config)
	logger_default.Close()
	logger_default = nil
	takeup = false
}


// 公共方法，defualt init
func defaultLoggerInit(lc *LogConfig) {
	if takeup == false {
		if lc == nil {
			// 设置默认操作！
			config = LogConfigDefalut
		} else {
			config = lc
		}
		logger_default = LoggerContainer(config)
	}
}