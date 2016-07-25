package libv2ray

import (
	"encoding/json"
	"errors"
	"os"

	"log"

	simplejson "github.com/bitly/go-simplejson"
)

type libv2rayconf struct {
	upscript      string
	downscript    string
	additionalEnv []string
	esco          []libv2rayconfEscortTarget
	rend          []libv2rayconfRenderCfgTarget
	vpnConfig     vpnserviceConfig
}

type libv2rayconfEscortTarget struct {
	Target      string   `json:"Target"`
	Args        []string `json:"Args"`
	Forgiveable bool     `json:"Forgiveable"`
}

type libv2rayconfRenderCfgTarget struct {
	Target string   `json:"Target"`
	Args   []string `json:"Args"`
	Source string   `json:"Source"`
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
	v.conf.additionalEnv, err = libconf.GetPath("env").StringArray()
	if err != nil {
		v.Callbacks.OnEmitStatus(-2, err.Error())
		return err
	}

	escortconfjson, exist := libconf.CheckGet("escort")
	if exist {
		/*
			var ok bool
			v.conf.esco, ok = escortconfjson.Interface().([]libv2rayconfEscortTarget)
			if !ok {
				v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config escort")
				return errors.New("Failed Type Assert: Config escort")
			}*/
		esco, ok := escortconfjson.MarshalJSON()
		if ok != nil {
			v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config escort")
			return errors.New("Failed Type Assert: Config escort")
		}
		err := json.Unmarshal(esco, &v.conf.esco)
		if err != nil {
			v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config escortX")
			return errors.New("Failed Type Assert: Config escortX")
		}

	}

	renderconfjson, exist := libconf.CheckGet("render")
	if exist {
		rend, ok := renderconfjson.MarshalJSON()
		if ok != nil {
			v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config render")
			return errors.New("Failed Type Assert: Config render")
		}
		err := json.Unmarshal(rend, &v.conf.rend)
		if err != nil {
			v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config renderX")
			log.Println(err, "Failed Type Assert: Config renderX")
			return errors.New("Failed Type Assert: Config renderX")
		}
	}

	vpnConfigconfjson, exist := libconf.CheckGet("vpnservice")
	if exist {
		vpnConfig, ok := vpnConfigconfjson.MarshalJSON()
		if ok != nil {
			v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config vpnConfig")
			return errors.New("Failed Type Assert: Config vpnConfig")
		}
		err := json.Unmarshal(vpnConfig, &v.conf.vpnConfig)
		if err != nil {
			v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config vpnConfigX")
			log.Println(err, "Failed Type Assert: Config vpnConfigX")
			return errors.New("Failed Type Assert: Config vpnConfigX")
		}
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
