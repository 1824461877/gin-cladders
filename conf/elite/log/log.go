package log

import (
	"errors"
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)


const tunnel_size_default = 1024

// Log config interface , Upper mouth
type LogConfigInter interface {
	Init() *LogConfig
}

// write file and console information, Lower mouth
// console log file a86	nd wirte log file interface
type WriteAll interface {
	Write(*ConsoleAndWrite,*Record) error
}

// write interface , write file this is it
type Write interface{
	Validator(*LogConfig) []FileWriteContainer // 正确日志验证器
	Init() (*FileFunc,error) // 实例化
	Write(*ConsoleAndWrite,*Record) error // 写入文件操作
	Flusher(*FileFunc) error // 销毁模式
}

// console interface , console log information
type ConsoleInfo interface {
	Validator(*LogConfig) bool // 文件输出接口
	SetColor(*LogConfig) *ConsoleFunc // 颜色设置接口
	String() string
}

// extend func ，have any extend implement
type AddFunc = func(config *LogConfig)

// LogWriteFile config ， defaults decide whether to use default parameters
// logwritefile is Write interface
// logconsolefile is consoleInnfo interface

// 写日子 和 输出日子一个总结
type SetUpLogWriteFile struct {
	defaults bool // 是不是默认设置
	logwritefile Write
	logconsolefile ConsoleInfo
	additional []AddFunc
}

// default func ， 默认配置模式
func LoggerContainer(lg LogConfigInter) *Logger {
	// 对已经实例化进行拦截操作!
	if logger_default != nil && takeup == false {
		takeup = true	//默认启动标志
		return logger_default
	}

	swf := &SetUpLogWriteFile{
		defaults:     true, // 配置默认模式
		logwritefile: &FileWriter{}, // 写日子
		logconsolefile: &colorRecord{}, // 读日子
	}
	// 调用 loggercontainer 集装箱，
	logger := swf.loggercontainer(lg)
	return logger
}

// custom func， 自定义配置模式
func CustomLoggerWriteFile(nf Write,cf ConsoleInfo,lg LogConfigInter,funcs ...AddFunc) *Logger {
	swf := &SetUpLogWriteFile{
		defaults:     false,
		logwritefile: nf,
		logconsolefile: cf,
		additional:funcs,
	}
	logger := swf.loggercontainer(lg)
	return logger
}

// log pipeline , implement Write and ConsoleInfo
// 执行logger server_container 的方法
func (swf *SetUpLogWriteFile) loggercontainer(lc LogConfigInter) *Logger {
	d := lc.Init()

	// default set func，默认是没有扩展方法
	if !swf.defaults {
		for _,v := range swf.additional {
			v(d)
		}
	}

	// logger struct
	logger := BoostrapLogWriter()

	// 输出日志操作！
	swf.consoleLogFile(logger,d)

	if logger.level = d.SetLeveL();  logger.level == -1 {
		return nil
	}
	if err := swf.writeLogFile(logger,d); err != nil {
		return nil
	}
	return logger
}

func (swf *SetUpLogWriteFile) consoleLogFile(logger *Logger,d *LogConfig) {
	if ok := swf.logconsolefile.Validator(d); !ok {
		return
	}
	color := swf.logconsolefile.SetColor(d)
	logger.RegisterConsoleInfo(color)
}

func (swf *SetUpLogWriteFile) writeLogFile(logger *Logger,d *LogConfig) error {
	valid := swf.logwritefile.Validator(d)
	if !valid[0].Bool {
		return errors.New("log is error")
	}

	if !valid[1].Bool  {
		logger.RegisterWriteFile(valid[0].Write,TRACE,FATAL)
	} else {
		logger.RegisterWriteFile(valid[0].Write,TRACE,INFO)
		logger.RegisterWriteFile(valid[1].Write,WARNING,ERROR)
	}
	return nil
}

// Execution Pool Struct
type WriteFuncs struct {
	writeslist WriteAll
	caw ConsoleAndWrite
}

// console struct and write struct  , gave file_console and write_console
type ConsoleAndWrite struct {
	filefuncs *FileFunc
	consolefunc *ConsoleFunc
}

type Logger struct {
	// Execution Pool
	writers []WriteFuncs

	// level grade , view level.go
	level int

	// time.now , unix time
	lastTime int64

	// production and consumption chan
	tunnel      chan *Record

	// temporal origin
	lastTimeStr string

	// error chan
	c chan bool

	// temporal origin
	layout string

	// sync.Pool , this is buffer pool. it need to be done
	recordPool *sync.Pool
}

// return Logger struct , instantiation Logger struct
func BoostrapLogWriter() *Logger {
	l := &Logger{
		writers:     []WriteFuncs{},
		level:       DEBUG,
		tunnel:make(chan *Record, tunnel_size_default),
		c:           make(chan bool,2),
		layout:      "2006/01/02 15:04:05",
		recordPool: &sync.Pool{New: func() interface{} {
			return &Record{}
		}},
	}

	// monitor tunnel
	go boostrapLogWriter(l)
	return l
}

// register consoleinfo , terminal display log information
func (l *Logger) RegisterConsoleInfo(color *ConsoleFunc) {
	wf := new(WriteFuncs)
	wf.writeslist = color
	wf.caw.consolefunc = nil
	l.writers = append(l.writers,*wf)
}

// register writefile , writer new log file
func (l *Logger) RegisterWriteFile(w Write,leve_floor,leve_ceil int) {
	filefunc , err := w.Init()
	if err != nil {
		panic(err)
	}
	filefunc.SetLogLevelFloor(leve_floor)
	filefunc.SetLogLevelCeil(leve_ceil)
	wf := new(WriteFuncs)
	wf.writeslist = w
	wf.caw.filefuncs = filefunc
	l.writers = append(l.writers,*wf)
}

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.writeslist.Write(&w.caw,r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	//rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}
			for _, w := range logger.writers {
				if err := w.writeslist.Write(&w.caw,r); err != nil {
					log.Println(err)
				}
			}

			logger.recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.writeslist.(Write); ok {
					if err := f.Flusher(w.caw.filefuncs); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		//case <-rotateTimer.C:
		//	for _, w := range logger.writers {
		//		if r, ok := w.(Rotater); ok {
		//			if err := r.Rotate(); err != nil {
		//				log.Println(err)
		//			}
		//		}
		//	}
		//	rotateTimer.Reset(time.Second * 10)
		}
	}
}

// production log information on the tunner chan
func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf, code string

	if level < l.level {
		fmt.Println("ok")
		return
	}


	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}
	//  获取缓存池的struct
	r := l.recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level

	l.tunnel <- r
}

// close tunnel chan
func (l *Logger) Close() {
	close(l.tunnel)
	<-l.c
	for _, w := range l.writers {
		if f, ok := w.writeslist.(Write); ok {
			if err := f.Flusher(w.caw.filefuncs); err != nil {
				log.Println(err)
			}
		}
	}
}

func (l *Logger) Trace(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(TRACE, fmt, args...)
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}
