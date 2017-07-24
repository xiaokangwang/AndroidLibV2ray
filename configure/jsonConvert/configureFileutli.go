package jsonConvert

import (
	"os"
)

type cfgtmpvars struct {
	fd *os.File
}

/*
func (v *jsonToPbConverter) parseCfg() io.Reader {
	//Use Context if possible

		if v.Context != nil {
			v.ConfigureFile, _ = v.Context.ReadProp(configureFile)
		}

	if v.ConfigureFile == "V2Ray_internal/ConfigureFileContent" {
		return strings.NewReader(v.ConfigureFileContent)
	}

	fd, err := os.Open(v.ConfigureFile)
	if err != nil {
		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")
	}
	v.cfgtmpvarsi.fd = fd
	v2rayJsonWithComment := &v2rayJsonWithComment.Reader{Reader: fd}
	return v2rayJsonWithComment
}

/*
func (v *jsonToPbConverter) parseCfgDone() {
	if v.cfgtmpvarsi.fd != nil {
		v.cfgtmpvarsi.fd.Close()
	}

}
*/
