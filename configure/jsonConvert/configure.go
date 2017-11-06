package jsonConvert

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"log"

	simplejson "github.com/bitly/go-simplejson"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
	v2rayJsonWithComment "v2ray.com/ext/encoding/json"
)

type libv2rayconf struct {
	upscript      string
	downscript    string
	additionalEnv []string
	esco          []libv2rayconfEscortTarget
	rend          []libv2rayconfRenderCfgTarget
	vpnConfig     vpnserviceConfig
	dnsloopfix    vpnserviceDnsloopFix
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
type JsonToPbConverter struct {
	conf    *libv2rayconf
	reading string
	Datadir string
	Cfgfile string
	Env     *configure.EnvironmentVar
}

func (v *JsonToPbConverter) parseConf() error {
	jsoncf, err := simplejson.NewFromReader(v.StripComment(v.reading))
	if err != nil {
		//v.Callbacks.OnEmitStatus(-2, err.Error())
		return err
	}
	libconf, isexist := jsoncf.CheckGet("#lib2ray")

	if !isexist {
		log.Print("No Vendor Conf found.")
		return nil
	}
	enabled, err := libconf.GetPath("enabled").Bool()
	if err != nil {
		//v.Callbacks.OnEmitStatus(-2, err.Error())
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
		//v.Callbacks.OnEmitStatus(-2, err.Error())
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
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config escort")
			return errors.New("Failed Type Assert: Config escort")
		}
		err := json.Unmarshal(esco, &v.conf.esco)
		if err != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config escortX")
			return errors.New("Failed Type Assert: Config escortX")
		}

	}

	renderconfjson, exist := libconf.CheckGet("render")
	if exist {
		rend, ok := renderconfjson.MarshalJSON()
		if ok != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config render")
			return errors.New("Failed Type Assert: Config render")
		}
		err := json.Unmarshal(rend, &v.conf.rend)
		if err != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config renderX")
			log.Println(err, "Failed Type Assert: Config renderX")
			return errors.New("Failed Type Assert: Config renderX")
		}
	}

	vpnConfigconfjson, exist := libconf.CheckGet("vpnservice")
	if exist {
		vpnConfig, ok := vpnConfigconfjson.MarshalJSON()
		if ok != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config vpnConfig")
			return errors.New("Failed Type Assert: Config vpnConfig")
		}
		err := json.Unmarshal(vpnConfig, &v.conf.vpnConfig)
		if err != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config vpnConfigX")
			log.Println(err, "Failed Type Assert: Config vpnConfigX")
			return errors.New("Failed Type Assert: Config vpnConfigX")
		}
	}

	vpndnsloopFix, exist := libconf.CheckGet("preparedDomainName")
	if exist {
		vpndnsloopFixJ, ok := vpndnsloopFix.MarshalJSON()
		if ok != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config vpndnsloopFixJ")
			return errors.New("Failed Type Assert: Config vpndnsloopFixJ")
		}
		err := json.Unmarshal(vpndnsloopFixJ, &v.conf.dnsloopfix)
		if err != nil {
			//v.Callbacks.OnEmitStatus(-2, "Failed Type Assert: Config vpndnsloopFixJ")
			log.Println(err, "Failed Type Assert: Config vpndnsloopFixJ")
			return errors.New("Failed Type Assert: Config vpndnsloopFixJ")
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

func (v *JsonToPbConverter) StripComment(Content string) io.Reader {
	Configure := strings.NewReader(Content)
	v2rayJsonWithComment := &v2rayJsonWithComment.Reader{Reader: Configure}
	var stp bytes.Buffer
	io.Copy(&stp, v2rayJsonWithComment)
	return &stp
}

func (v *JsonToPbConverter) LoadFromString(ctx string) {
	v.reading = ctx
}
func (v *JsonToPbConverter) LoadFromFile(loc string) error {
	file, err := ioutil.ReadFile(loc)
	if err != nil {
		return err
	}
	v.reading = string(file)
	return nil
}

func (v *JsonToPbConverter) Parse() error {
	err := v.parseConf()
	if err != nil {
		return err
	}
	//ConvertEnv
	v.Env = &configure.EnvironmentVar{}
	if v.conf == nil {
		return nil
	}
	v.Env.Vars = envToMap(v.conf.additionalEnv)
	v.renderAll()
	return nil
}

func (v *JsonToPbConverter) ToPb() *configure.LibV2RayConf {
	if v.conf == nil {
		return nil
	}
	return ConvertToPb(*v.conf)
}
