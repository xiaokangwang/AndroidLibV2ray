package Process

import (
	"fmt"
	"os"
	"strconv"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
)

func (v *EnvironmentCreater) AddEnvironment(env []string) []string {
	datadir := v.Context.GetDataDir()
	env = append(env, "proxyuid="+strconv.Itoa(os.Getuid()))
	env = append(env, "datadir="+datadir)
	//env = append(env, "cfgdir="+filepath.Dir(v.Context.GetConfigureFile())+"")
	return env
}

func (v *EnvironmentCreater) GetEnvironment() []string {
	if v.Conf == nil {
		return v.AddEnvironment(make([]string, 0))
	}
	env := EnvJoins(v.Conf.Vars)
	env = v.AddEnvironment(env)
	return env
}

type EnvironmentCreater struct {
	Conf    *configure.EnvironmentVar
	Context *CoreI.Status
}

func EnvJoins(env map[string]string) []string {
	ret := make([]string, 0, len(env))
	for i, v := range env {
		ret = append(ret, fmt.Sprintf("%v=%v", i, v))
	}
	return ret
}
