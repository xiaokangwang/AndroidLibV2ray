package libv2ray

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
)

func NewLib2rayContext() *V2RayContext {
	return new(V2RayContext)
}

type V2RayContext struct {
	configureFile string
	Callbacks     V2RayContextCallbacks
	PackageName   string
	Status        CoreI.Status
}

const configureFile = "ConfigureFile"

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
	none := func() *StringArrayList {
		var retsg []string
		retsg = append(retsg, "..")
		return &StringArrayList{list: retsg}
	}
	if vc.GetConfigureFile() == "" {
		return none()
	}
	dir := path.Dir(vc.configureFile)
	dfd, err := os.Open(dir)
	if err != nil {
		return none()
	}
	d, err := dfd.Readdirnames(128)
	if err != nil {
		return none()
	}
	d = append(d, "..")
	for di := range d {
		d[di] = path.Dir(vc.GetConfigureFile()) + "/" + d[di]
	}
	return &StringArrayList{list: d}
}

func (vc *V2RayContext) GetBriefDesc(pathn string) string {
	_, ret := path.Split(pathn)
	return ret
}

func (vc *V2RayContext) AssignConfigureFile(cf string) {
	if strings.HasSuffix(cf, "..") {
		vc.Callbacks.OnRefreshNeeded()
		vc.Callbacks.OnFileSelectTriggerd()
		return
	}
	log.Print(cf)
	vc.configureFile = cf
	vc.WriteProp(configureFile, cf)
	vc.Callbacks.OnRefreshNeeded()
}

func (vc *V2RayContext) GetConfigureFile() string {
	if vc.configureFile == "" {
		vc.configureFile, _ = vc.ReadProp(configureFile)
	}
	return vc.configureFile
}

type V2RayContextCallbacks interface {
	OnRefreshNeeded()
	OnFileSelectTriggerd()
}

func (vc *V2RayContext) ReadProp(name string) (string, error) {
	os.MkdirAll(vc.Status.GetDataDir()+"config", 0700)
	fd, err := os.Open(vc.Status.GetDataDir() + "config/" + name)
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

func (vc *V2RayContext) ReadPropD(name string) string {
	ctx, _ := vc.ReadProp(name)
	return ctx
}

func (vc *V2RayContext) WriteProp(name string, cont string) error {
	os.MkdirAll(vc.Status.GetDataDir()+"config", 0700)
	return ioutil.WriteFile(vc.Status.GetDataDir()+"config/"+name, []byte(cont), 0600)
}
