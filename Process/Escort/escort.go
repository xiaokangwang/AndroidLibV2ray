package Escort

import (
	"os"
	"os/exec"

	"log"
)
import "github.com/xiaokangwang/AndroidLibV2ray/configure"
import "github.com/xiaokangwang/AndroidLibV2ray/CoreI"
import "github.com/xiaokangwang/AndroidLibV2ray/Process"

func (v *Escorting) EscortRun(proc string, pt []string, forgiveable bool, tapfd int) {
	count := 42
	for count > 0 {
		cmd := exec.Command(proc, pt...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		ect := Process.EnvironmentCreater{Conf: v.Env, Context: v.status}
		env := ect.GetEnvironment()
		env = append(env, os.Environ()...)
		env = ect.AddEnvironment(env)
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
		if v.status.IsRunning {
			log.Println("Unexpected Exit")
			count--
		} else {
			return
		}
	}

	if v.status.IsRunning && !forgiveable {
		v.unforgivnesschan <- 0
	}

}

func (v *Escorting) escortBeg(proc string, pt []string, forgiveable bool) {
	go v.EscortRun(proc, pt, forgiveable, 0)
}

func (v *Escorting) unforgivenessCloser() {
	log.Println("unforgivenessCloser() <-v.unforgivnesschan")
	<-v.unforgivnesschan
	/*if v.status.IsRunning {
		//TODO:v.caller.StopLoop()
		log.Println("Closed As unforgivenessCloser decided so.")

	}*/
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

func (v *Escorting) EscortingUP() {
	if v.escortProcess != nil {
		return
	}
	v.escortProcess = new([](*os.Process))
	go v.unforgivenessCloser()
	for _, esct := range v.Configure {
		v.escortBeg(esct.Target, esct.Args, esct.Forgiveable)
	}
}
func (v *Escorting) EscortingUPV() {
	if v.escortProcess != nil {
		return
	}
	v.escortProcess = new([](*os.Process))
	go v.unforgivenessCloser()
}
func (v *Escorting) EscortingDown() {
	log.Println("escortingDown() Killing all escorted process ")
	if v.escortProcess == nil {
		return
	}
	for _, pr := range *v.escortProcess {
		pr.Kill()
	}
	log.Println("escortingDown() v.unforgivnesschan <- 0")
	select {
	case v.unforgivnesschan <- 0:
	}
	v.escortProcess = nil

}

func (v *Escorting) SetStatus(st *CoreI.Status) {
	v.status = st
}

func NewEscort() *Escorting {
	return &Escorting{unforgivnesschan: make(chan int)}
}

type Escorting struct {
	escortProcess    *[](*os.Process)
	unforgivnesschan chan int
	status           *CoreI.Status
	Configure        []*configure.EscortedProcess
	Env              *configure.EnvironmentVar
}
