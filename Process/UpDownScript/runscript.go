package UpDownScript

import (
	"os"
	"os/exec"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
)

func (v *UpDownScript) Runbash(cc string, env []string) error {
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
	Configure *configure.UpDownScripts
	Env       *configure.EnvironmentVar
}

func (v *UpDownScript) RunUpScript() {
	bashs := v.Configure.UpScript
	ect := Process.EnvironmentCreater{Conf: v.Env, Context: v.status}
	env := ect.GetEnvironment()
	env = append(env, os.Environ()...)
	env = ect.AddEnvironment(env)
	v.Runbash(bashs, env)
}
func (v *UpDownScript) RunDownScript() {
	bashs := v.Configure.DownScript
	ect := Process.EnvironmentCreater{Conf: v.Env, Context: v.status}
	env := ect.GetEnvironment()
	env = append(env, os.Environ()...)
	env = ect.AddEnvironment(env)
	v.Runbash(bashs, env)
}

func (v *UpDownScript) SetStatus(st *CoreI.Status) {
	v.status = st
}
