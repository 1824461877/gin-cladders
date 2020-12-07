package log

type colorRecord = Record

func (r *colorRecord) Validator(logger *LogConfig) bool {
	if logger.CW.On {
		return true
	}
	return false
}

func (r *colorRecord) SetColor(logger *LogConfig) *ConsoleFunc {
	if logger.CW.Color {
		return &ConsoleFunc{true}
	}
	return &ConsoleFunc{false}
}

func (r *colorRecord) OutputTemplate() string {
	return r.TemplateOutput()
}

