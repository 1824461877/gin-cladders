package log

import "fmt"

var (
	LEVEL = []string{
		"TRACE",
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}
)


const (
	TRACE = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

type Record struct {
	time string
	code string
	info string
	level int
}

func (r Record) String() string{
	return fmt.Sprintf("[%s][%s][%s] %s\n", LEVEL[r.level], r.time, r.code, r.info)
}

func (r Record) TemplateOutput() string {
	switch r.level {
	case TRACE:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL[r.level], r.code, r.info)

	case DEBUG:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[34m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL[r.level], r.code, r.info)

	case INFO:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[32m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL[r.level], r.code, r.info)

	case WARNING:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[33m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL[r.level], r.code, r.info)

	case ERROR:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[31m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL[r.level], r.code, r.info)

	case FATAL:
		return fmt.Sprintf("\033[36m%s\033[0m [\033[35m%s\033[0m] \033[47;30m%s\033[0m %s\n",
			r.time, LEVEL[r.level], r.code, r.info)
	default:
		return fmt.Sprint("no level")
	}
}
