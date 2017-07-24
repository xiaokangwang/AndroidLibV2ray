package libv2ray

import (
	"os"
	"os/exec"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
)

func (v *UpDownScript) runbash(cc string, env []string) error {
	cmd := exec.Command("/system/bin/sh", "-c", cc)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ect := Process.EnvironmentCreater{Conf: v.Env, Context: v.status}
	env = append(env, os.Environ()...)
	env = ect.AddEnvironment(env)
	cmd.Env = env
	err := cmd.Run()
	return err
}

type UpDownScript struct {
	status    *CoreI.Status
	configure *configure.UpDownScripts
	Env       *configure.EnvironmentVar
}
