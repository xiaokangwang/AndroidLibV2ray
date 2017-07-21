package libv2ray

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (v *jsonToPbConverter) addEnvironment(env []string) []string {
	datadir := v.Datadir
	env = append(env, "proxyuid="+strconv.Itoa(os.Getuid()))
	env = append(env, "datadir="+datadir)
	env = append(env, "cfgdir="+filepath.Dir(v.Cfgfile)+"")
	return env
}

func (v *jsonToPbConverter) getEnvironment() []string {
	if v.conf == nil {
		return v.addEnvironment(make([]string, 0))
	}
	env := v.conf.additionalEnv
	env = v.addEnvironment(env)
	return env
}

func envToMap(k []string) map[string]string {
	var themap map[string]string
	themap = make(map[string]string)
	for _, val := range k {
		r := strings.Index(val, "=")
		themap[val[:r]] = val[r+1:]
	}
	return themap
}
