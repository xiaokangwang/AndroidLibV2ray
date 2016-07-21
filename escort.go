package libv2ray

import (
	"os"
	"os/exec"

	"log"
)

func (v *V2RayPoint) escortRun(proc string, pt []string, forgiveable bool, tapfd int) {
	count := 42
	for count > 0 {
		cmd := exec.Command(proc, pt...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		env := v.conf.additionalEnv
		env = append(env, os.Environ()...)
		env = v.addEnvironment(env)
		cmd.Env = env

		if tapfd != 0 {
			file := os.NewFile(uintptr(tapfd), "/dev/tap0")
			var files []*os.File
			cmd.ExtraFiles = append(files, file)
		}

		err := cmd.Start()
		log.Println(proc)
		log.Println(pt)
		if err != nil {
			log.Println(err)
		}
		*v.escortProcess = append(*v.escortProcess, cmd.Process)
		log.Println("Waiting....")
		err = cmd.Wait()
		log.Println("Exit")
		log.Println(err)
		if v.IsRunning {
			log.Println("Unexpected Exit")
			count--
		} else {
			return
		}
	}

	if v.IsRunning && !forgiveable {
		v.unforgivnesschan <- 0
	}

}

func (v *V2RayPoint) escortBeg(proc string, pt []string, forgiveable bool) {
	go v.escortRun(proc, pt, forgiveable, 0)
}

func (v *V2RayPoint) unforgivenessCloser() {
	log.Println("unforgivenessCloser() <-v.unforgivnesschan")
	<-v.unforgivnesschan
	if v.IsRunning {
		v.StopLoop()
		log.Println("Closed As unforgivenessCloser decided so.")
		v.Callbacks.OnEmitStatus(0, "Closed As unforgivenessCloser decided so.")
	}
	remain := true
	for remain {
		select {
		case <-v.unforgivnesschan:
			log.Println("unforgivenessCloser() removing reminder unforgivness sign")
			break
		default:
			remain = false
		}
	}
	log.Println("unforgivenessCloser() quit")
}

func (v *V2RayPoint) escortingUP() {
	if v.escortProcess != nil {
		return
	}
	v.escortProcess = new([](*os.Process))
	go v.unforgivenessCloser()
	for _, esct := range v.conf.esco {
		v.escortBeg(esct.Target, esct.Args, esct.Forgiveable)
	}
}
func (v *V2RayPoint) escortingDown() {
	log.Println("escortingDown() Killing all escorted process ")
	for _, pr := range *v.escortProcess {
		pr.Kill()
	}
	log.Println("escortingDown() v.unforgivnesschan <- 0")
	select {
	case v.unforgivnesschan <- 0:
	}
	v.escortProcess = nil

}
