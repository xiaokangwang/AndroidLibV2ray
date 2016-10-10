package libv2ray

import (
	"io"
	"os"
	"strings"
)

type cfgtmpvars struct {
	fd *os.File
}

func (v *V2RayPoint) parseCfg() io.Reader {
	if v.ConfigureFile == "V2Ray_internal/ConfigureFileContent" {
		return strings.NewReader(v.ConfigureFileContent)
	}

	fd, err := os.Open(v.ConfigureFile)
	if err != nil {
		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")
	}
	v.cfgtmpvarsi.fd = fd
	return fd
}
func (v *V2RayPoint) parseCfgDone() {
	if v.cfgtmpvarsi.fd != nil {
		v.cfgtmpvarsi.fd.Close()
	}

}
