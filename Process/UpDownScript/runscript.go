package libv2ray

import (
	"os"
	"os/exec"
)

func (v *V2RayPoint) runbash(cc string, env []string) error {
	cmd := exec.Command("/system/bin/sh", "-c", cc)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env = append(env, os.Environ()...)
	env = v.addEnvironment(env)
	cmd.Env = env
	err := cmd.Run()
	return err
}
