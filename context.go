package libv2ray

import (
	"io/ioutil"
	"os"
	"path"
)

func NewLib2rayContext() *V2RayContext {
	return new(V2RayContext)
}

type V2RayContext struct {
	configureFile string
	Callbacks     *V2RayContextCallbacks
	PackageName   string
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

func (vc *V2RayContext) ListConfigureFileDir() *StringArrayList {
	dir := path.Dir(vc.configureFile)
	dfd, err := os.Open(dir)
	if err != nil {
		return nil
	}
	d, err := dfd.Readdirnames(128)
	if err != nil {
		return nil
	}
	return &StringArrayList{list: d}
}

func (vc *V2RayContext) GetBriefDesc(pathn string) string {
	_, ret := path.Split(pathn)
	return ret
}

func (vc *V2RayContext) AssignConfigureFile(cf string) {
	vc.configureFile = cf
}

func (vc *V2RayContext) GetConfigureFile() string {
	return vc.configureFile
}

type V2RayContextCallbacks interface {
	OnRefreshNeeded()
}

func (vc *V2RayContext) ReadProp(name string) (string, error) {
	os.MkdirAll(vc.getDataDir()+"config", 0700)
	fd, err := os.Open(vc.getDataDir() + "config/" + name)
	if err != nil {
		return "", err
	}
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return "", err
	}
	fd.Close()
	return string(content), nil
}

func (vc *V2RayContext) WriteProp(name string, cont string) error {
	os.MkdirAll(vc.getDataDir()+"config", 0700)
	return ioutil.WriteFile(vc.getDataDir()+"config/"+name, []byte(cont), 0600)
}

func (v *V2RayContext) getDataDir() string {
	var datadir = "/data/data/org.kkdev.v2raygo/"
	if v.PackageName != "" {
		datadir = "/data/data/" + v.PackageName + "/"
	}
	return datadir
}
