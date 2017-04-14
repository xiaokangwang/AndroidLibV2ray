package libv2ray

import (
	"os"
	"path"
)

func NewLib2rayContext() *V2RayContext {
	return new(V2RayContext)
}

type V2RayContext struct {
	configureFile string
}

func (vc *V2RayContext) CheckConfigureFile() bool {
	//Check if file exist
	if !exists(vc.configureFile) {
		return false
	}
	return true
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func (vc *V2RayContext) ListConfigureFileDir() []string {
	dir := path.Dir(vc.configureFile)
	dfd, err := os.Open(dir)
	if err != nil {
		return nil
	}
	d, err := dfd.Readdirnames(128)
	if err != nil {
		return nil
	}
	return d
}
func (vc *V2RayContext) AssignConfigureFile(cf string) {
	vc.configureFile = cf
}

func (vc *V2RayContext) GetConfigureFile() string {
	return vc.configureFile
}
