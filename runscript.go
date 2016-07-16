package libv2ray

import "os/exec"

func runbash(cc string, env []string) error {
	cmd := exec.Command("/system/bin/sh", "-c", cc)
	cmd.Env = env
	err := cmd.Run()
	return err
}
