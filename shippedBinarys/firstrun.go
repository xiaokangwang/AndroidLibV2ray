package libv2ray

import (
	"errors"
	"io"
	"os"
	"runtime"
	"strconv"

	"golang.org/x/mobile/asset"
)

func (v *V2RayPoint) checkIfRcExist() error {
	datadir := v.getDataDir()

	if _, err := os.Stat(datadir + strconv.Itoa(CheckVersion())); !os.IsNotExist(err) {
		return nil
	}
	var arcIndepRc []string
	arcIndepRc = append(arcIndepRc, "pdnsd-te.conf")
	var arcDepRc []string
	arcDepRc = append(arcDepRc, "pdnsd", "tun2socks")
	for _, rcn := range arcIndepRc {
		f, err := asset.Open(rcn)
		if err != nil {
			return err
		}
		fw, err := os.OpenFile(datadir+rcn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, f)
		if err != nil {
			return err
		}
		f.Close()
		fw.Close()
	}
	var sarch string
	switch runtime.GOARCH {
	case "amd64":
		sarch = "x86_64"
		break
	case "386":
		sarch = "x86"
		break
	case "arm":
		sarch = "armeabi-v7a"
		break
	case "arm64":
		sarch = "arm64-v8a"
		break
	default:
		return errors.New("Unsupported Arch")
	}

	for _, rcn := range arcDepRc {
		f, err := asset.Open(sarch + "/" + rcn)
		if err != nil {
			return err
		}
		fw, err := os.OpenFile(datadir+rcn, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, f)
		if err != nil {
			return err
		}
		f.Close()
		fw.Close()
	}

	s, _ := os.Create(datadir + strconv.Itoa(CheckVersion()))
	s.Close()

	return nil
}

func (v *V2RayPoint) getDataDir() string {
	var datadir = "/data/data/org.kkdev.v2raygo/"
	if v.PackageName != "" {
		datadir = "/data/data/" + v.PackageName + "/"
	}
	return datadir
}
