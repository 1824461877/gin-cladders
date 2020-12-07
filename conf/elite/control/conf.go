package control

import (
	dlog "gin-cladder/conf/elite/log"
)

// 基本配置
type BaseConf struct {
	DebugMode string `mapstructure:"debug_mode"`  // debug 模式选型
	TimeLocation string `mapstructure:"time_location"` // 时区
	Log LogConfig `mapstructure:"log"` // 日志继承
	Http HttpConfig `mapstructure:"http"`
	Base         struct {
		DebugMode    string `mapstructure:"debug_mode"`
		TimeLocation string `mapstructure:"time_location"`
	} `mapstructure:"base"`
}


type HttpConfig struct {
	Addr string `mapstructure:"addr"`
	ReadTimeout int `mapstructure:"read_timeout"`
	WriteTimeput int `mapstructure:"write_timeout"`
	MaxHeaderBytes int `mapstructure:"max_header_bytes"`
	AllowHost   []string `mapstructure:"allow_host"`
}

// 日志配置
type LogConfig struct {
	LogLevel string `mapstructure:"log_level"`
	FW ConfWriterFile `mapstructure:"file_writer"`
	CW ConfConsoleFile `mapstructure:"console_writer"`
}

// 日志文件写入模块
type ConfWriterFile struct {
	On bool `mapstructure:"on"`
	LogPath string `mapstructure:"log_path"`
	WFLogPath string `mapstructure:"wf_log_path"`
}

// 日志文件输出模块
type ConfConsoleFile struct {
	On bool `mapstructure:"on"`
	Color bool `mapstructure:"color"`
}


// 实例化 BaseConf
func InitBaseConf(path string) error {
	// 实例一个 baseconf
	ConfBase = &BaseConf{}

	// 解析文件中的配置操作！把解析的结果放入 ConfBase 里面
	if err := ParseConfig(path,ConfBase); err != nil {
		return err
	}
	// 配置文件模式
	if ConfBase.DebugMode == "" {
		if ConfBase.Base.DebugMode != "" {
			ConfBase.DebugMode = ConfBase.Base.DebugMode
		} else {
			ConfBase.DebugMode = "debug"
		}
	}
	if ConfBase.TimeLocation == "" {
		if ConfBase.Base.TimeLocation != "" {
			ConfBase.TimeLocation = ConfBase.Base.TimeLocation
		} else {
			ConfBase.TimeLocation = "Asia/Chongqing"
		}
	}
	if ConfBase.Log.LogLevel == "" {
		ConfBase.Log.LogLevel = "trace"
	}

	// 日子基本配置
	s := &dlog.LogConfig{
		Level: ConfBase.Log.LogLevel, // 日志的模式
		FW:    dlog.ConfFileWriter{ // 后台写入文件操作
			On:        ConfBase.Log.FW.On, // 是否开启
			LogPath:   ConfBase.Log.FW.LogPath, // 日志存放路径
			WFLogPath: ConfBase.Log.FW.WFLogPath, // 错误日志存放路径
		},
		CW:    dlog.ConfFileConsole{ // 前台对日志进行输出
			On:    ConfBase.Log.CW.On, // 是否开启
			Color: ConfBase.Log.CW.Color,  // 是否有颜色变化
		},
	}

	// default 日志
	dlog.SetupDefaultLogWithConf(s)
	return nil
}

