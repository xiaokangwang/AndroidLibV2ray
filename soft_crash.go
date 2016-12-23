package libv2ray

// Debug code removed

const _ = `
import (
	"io/ioutil"
	"log"
	"os"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"
)

/*
#include <unistd.h>
#include <sys/types.h>
#include <pwd.h>
#include <stdlib.h>
*/
import "C"

func (v *V2RayPoint) setupSoftCrashMonitor() {
	if !v.softcrashMonitor {
		v.softcrashMonitor = true
		go v.softCrashMonitor()
	}
}

func (v *V2RayPoint) softCrashMonitor() {
	currentUid := os.Getpid()
	currentUidS := strconv.Itoa(currentUid)
	getCpuStat := func(uid string) int {
		contents, err := ioutil.ReadFile("/proc/" + uid + "/stat")
		if err != nil {
			log.Fatal("CrashMonitor is not working,1")
			return 0
		}
		ps := strings.Split(string(contents), " ")
		time, _ := strconv.Atoi(ps[13])
		return time
	}
	//Honey moon
	<-time.Tick(time.Second * 3)
	var sc_clk_tck C.long
	sc_clk_tck = C.sysconf(C._SC_CLK_TCK)
	last_usage := getCpuStat(currentUidS)

	capt := func() {
		log.Println("Debug: capturing")
		//debug high CPU usage
		file, _ := os.Create(v.getDataDir() + "goroutine_profile." + strconv.Itoa(int(time.Now().Unix())) + ".please_submit.bugreport")
		pprof.Lookup("goroutine").WriteTo(file, 1)
		log.Println("Debug: captured")
	}
	for {
		<-time.Tick(time.Second * 6)
		currentUsage := getCpuStat(currentUidS)
		log.Println("Using CPU ", currentUsage-last_usage, "/", int(sc_clk_tck))
		if currentUsage-last_usage >= (5 * int(sc_clk_tck)) {
			log.Println("CPU Usage TOO high")
			if v.isDebugTriggered() {
				capt()
			}
			log.Println("Crashing: CPU Usage too high.")
			v.Callbacks.OnEmitStatus(-3322, "SoftCrash")
			os.Exit(-100)
		}

		last_usage = currentUsage
	}
}
`
