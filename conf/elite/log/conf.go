package log

// 配置文件日志写入
type ConfFileWriter struct {
	On bool `toml:"On"`
	LogPath string `toml:"log_path"`
	WFLogPath string `toml:"wf_log_path"`
}

// 配置文件输出
type ConfFileConsole struct {
	On bool `toml:"On"`
	Color bool `toml:"Color"`
}

// 日志最基本模块
type LogConfig struct {
	Level string
	matching bool // 配置级别
	FW ConfFileWriter
	CW ConfFileConsole
}

// 默认 logconfig
var LogConfigDefalut = &LogConfig{
	Level: "debug",
	FW: ConfFileWriter{
		On:true,
		LogPath:"./log_test.log",
		WFLogPath:"./log_test.wf.log",
	},
	CW: ConfFileConsole{
		On:true,
		Color:true,
	},
}


func NewLogConfig() *LogConfig{
	return LogConfigDefalut
}

func (lg *LogConfig) Init() *LogConfig {
	return lg
}

func (lc *LogConfig) matchingLevel() bool{
	lc.matching = false
	for _,v := range LEVEL {
		if lc.Level == v {
			lc.matching = true
		}
	}
	return lc.matching
}

func (fc *LogConfig) SetLeveL() int {
	switch fc.Level {
	case "trace":
		return TRACE
	case "debug":
		return DEBUG
	case "info":
		return INFO
	case "warning":
		return WARNING
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return -1
	}
}

func SetupDefaultLogWithConf(lc *LogConfig) {
	// 实例一个 日志操作！
	defaultLoggerInit(lc)
}

func (fw *ConfFileWriter) GetLogPath() string {
	if fw.On {
		if fw.LogPath != "" {
			return fw.LogPath
		}
	}
	return ""
}

func (fw *ConfFileWriter) GetWFLogPath() string {
	if fw.On {
		if fw.WFLogPath != "" {
			return fw.WFLogPath
		}
	}
	return ""
}

func (fc *ConfFileConsole) Colors() bool {
	if fc.On {
		return fc.On
	}
	return false
}
