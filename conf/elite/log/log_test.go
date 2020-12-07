package log

import (
	"testing"
)

func TestLogInstance(t *testing.T) {
	s := &LogConfig{
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
	loggers := LoggerContainer(s)
	loggers.Info("this is trace")
	loggers.Close()
}
