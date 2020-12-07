package log


type FileWriter struct {
	filename      string
}

func (c *FileWriter) Validator(inter *LogConfig) []FileWriteContainer {
	var lopath,wflogpath = new(FileWriter),new(FileWriter)
	if inter.FW.On {
		if len(inter.FW.LogPath) > 0 {
			lopath.filename = inter.FW.GetLogPath()
			if lopath.filename == "" || inter.FW.GetWFLogPath() == "" {
				return nil
			}
			wflogpath.filename = inter.FW.GetWFLogPath()
		}
	}
	return []FileWriteContainer{{lopath,true},{wflogpath,true}}
}

func (c *FileWriter) Init() (*FileFunc,error) {
	filefunc , err := CreateLogFile(c.filename)
	if err != nil{
		return nil,err
	}
	return filefunc,nil
}

func (c *FileWriter) Write(caw *ConsoleAndWrite,r *Record) error {
	err := WriteFileBufWrite(caw.filefuncs,r)
	if err != nil {
		return err
	}
	return nil
}

func (c *FileWriter) Flusher(filefunc *FileFunc) error {
	err := Flush(filefunc)
	if err != nil {
		return err
	}
	return nil
}