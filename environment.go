package libv2ray

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (v *V2RayPoint) addEnvironment(env []string) []string {
	datadir := v.getDataDir()
	env = append(env, "proxyuid="+strconv.Itoa(os.Getuid()))
	env = append(env, "datadir="+datadir)
	env = append(env, "cfgdir="+filepath.Dir(v.ConfigureFile)+"")
	return env
}

func (v *V2RayPoint) getEnvironment() []string {
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
