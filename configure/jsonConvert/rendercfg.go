package jsonConvert

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process"
)

func (v *JsonToPbConverter) renderAll() {
	v.renderesco()
	v.renderptm()
	v.rendervpn()
}

func (v *JsonToPbConverter) renderptm() {
	Eg := &Process.EnvironmentCreater{Conf: v.Env, Context: &CoreI.Status{PackageName: v.Datadir}}
	envr0 := envToMap(Eg.GetEnvironment())
	mf0 := func(lookup string) string {
		if envl0, ok := envr0[lookup]; ok {
			return envl0
		}
		return ""
	}
	for _, r := range v.conf.rend {
		for key2 := range r.Args {
			r.Args[key2] = os.Expand(r.Args[key2], mf0)
		}
		mf := func(lookup string) string {
			envr := envToMap(append(Eg.GetEnvironment(), r.Args...))
			if envl, ok := envr[lookup]; ok {
				return envl
			}
			return ""
		}
		fs := os.Expand(r.Source, mf)
		ft := os.Expand(r.Target, mf)
		fds, err := os.Open(fs)
		if err != nil {
			log.Println(err, fs)
		}
		input, err := ioutil.ReadAll(fds)
		if err != nil {
			log.Println(err, fs, "RA")
		}
		inputs := string(input)
		op := os.Expand(inputs, mf)
		opn := strings.NewReader(op)

		fdt, err := os.Create(ft)
		if err != nil {
			log.Println(err, ft)
		}
		_, err = io.Copy(fdt, opn)
		if err != nil {
			log.Println(err, ft, "CP")
		}
	}
}
func (v *JsonToPbConverter) renderesco() {
	Eg := &Process.EnvironmentCreater{Conf: v.Env, Context: &CoreI.Status{PackageName: v.Datadir}}
	envr := envToMap(Eg.GetEnvironment())
	for key := range v.conf.esco {
		mf := func(lookup string) string {
			if envl, ok := envr[lookup]; ok {
				return envl
			}
			return ""
		}
		v.conf.esco[key].Target = os.Expand(v.conf.esco[key].Target, mf)
		for key2 := range v.conf.esco[key].Args {
			v.conf.esco[key].Args[key2] = os.Expand(v.conf.esco[key].Args[key2], mf)
		}
	}
}

func (v *JsonToPbConverter) rendervpn() {
	Eg := &Process.EnvironmentCreater{Conf: v.Env, Context: &CoreI.Status{PackageName: v.Datadir}}
	envr := envToMap(Eg.GetEnvironment())
	mf := func(lookup string) string {
		if envl, ok := envr[lookup]; ok {
			return envl
		}
		return ""
	}
	for key := range v.conf.vpnConfig.Args {
		v.conf.vpnConfig.Args[key] = os.Expand(v.conf.vpnConfig.Args[key], mf)
	}
	v.conf.vpnConfig.Target = os.Expand(v.conf.vpnConfig.Target, mf)
}
