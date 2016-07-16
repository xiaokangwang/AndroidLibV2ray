package libv2ray

import (
	"os"

	"log"

	simplejson "github.com/bitly/go-simplejson"
)

type libv2rayconf struct {
	upscript      string
	downscript    string
	additionalEnv []string
}

func (v *V2RayPoint) parseConf() error {
	fconffd, err := os.Open(v.ConfigureFile)
	if err != nil {
		v.Callbacks.OnEmitStatus(-2, "Failed to read config file ("+v.ConfigureFile+"):"+err.Error())
		return err
	}
	defer fconffd.Close()
	jsoncf, err := simplejson.NewFromReader(fconffd)
	if err != nil {
		v.Callbacks.OnEmitStatus(-2, err.Error())
		return err
	}
	libconf, isexist := jsoncf.CheckGet("#lib2ray")

	if !isexist {
		log.Print("No Vendor Conf found.")
		return nil
	}
	enabled, err := libconf.GetPath("enabled").Bool()
	if err != nil {
		v.Callbacks.OnEmitStatus(-2, err.Error())
		return err
	}
	if !enabled {
		return nil
	}

	v.conf = &libv2rayconf{}

	v.conf.upscript = jsonStringDefault(libconf.GetPath("listener", "onUp"), "#none")
	v.conf.downscript = jsonStringDefault(libconf.GetPath("listener", "onDown"), "#none")
	v.conf.additionalEnv, err = libconf.GetPath("listener", "env").StringArray()
	if err != nil {
		v.Callbacks.OnEmitStatus(-2, err.Error())
		return err
	}

	return nil
}

func jsonStringDefault(jsons *simplejson.Json, def string) string {
	s, err := jsons.String()
	if err != nil {
		return def
	}
	return s
}
