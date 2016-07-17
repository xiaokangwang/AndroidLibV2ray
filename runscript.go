package libv2ray

import (
	"os"
	"os/exec"
	"strconv"
)

func runbash(cc string, env []string) error {
	cmd := exec.Command("/system/bin/sh", "-c", cc)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	env = append(env, os.Environ()...)
	env = append(env, "proxyuid="+strconv.Itoa(os.Getuid()))
	env = append(env, "datadir="+datadir)
	cmd.Env = env
	err := cmd.Run()
	return err
}
