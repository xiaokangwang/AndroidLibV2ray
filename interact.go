package libv2ray

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/log"
	"v2ray.com/core/common/errors"

	v2rayconf "v2ray.com/ext/tools/conf/serial"
)
import "github.com/xiaokangwang/AndroidLibV2ray/configure"

/*V2RayPoint V2Ray Point Server
This is territory of Go, so no getter and setters!

Notice:
ConfigureFile can be either the path of config file or
"V2Ray_internal/ConfigureFileContent" in case you wish to

*/
type V2RayPoint struct {
	ConfigureFile        string
	ConfigureFileContent string
	Callbacks            V2RayCallbacks
	vpoint               core.Server
	IsRunning            bool
	//conf                 *libv2rayconf
	confng           *configure.LibV2RayConf
	escortProcess    *[](*os.Process)
	unforgivnesschan chan int
	VpnSupportSet    V2RayVPNServiceSupportsSet
	VpnSupportnodup  bool
	PackageName      string
	cfgtmpvarsi      cfgtmpvars
	//softcrashMonitor     bool
	prepareddomain  preparedDomain
	v2rayOP         *sync.Mutex
	interuptDeferto int64
	Context         *V2RayContext
}

/*V2RayCallbacks a Callback set for V2Ray
 */
type V2RayCallbacks interface {
	OnEmitStatus(int, string) int
}

func (v *V2RayPoint) pointloop() {
	//v.setupSoftCrashMonitor()

	v.VpnSupportnodup = false

	if v.parseConf() != nil {
		return
	}

	err := v.checkIfRcExist()

	if err != nil {
		log.Trace(errors.New("Failed to copy asset").Base(err).AtError())
		v.Callbacks.OnEmitStatus(-1, "Failed to copy asset ("+err.Error()+")")
	}

	log.Trace(errors.New("v.renderAll()"))
	v.renderAll()

	config, err := v2rayconf.LoadJSONConfig(v.parseCfg())
	if err != nil {
		log.Trace(errors.New("Failed to read config file (", v.ConfigureFile, "): ", v.ConfigureFile).Base(err).AtError())

		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")

		return
	}

	vPoint, err := core.New(config)
	if err != nil {
		log.Trace(errors.New("Failed to create Point server").Base(err))

		v.Callbacks.OnEmitStatus(-1, "Failed to create Point server ("+err.Error()+")")

		return
	}
	v.IsRunning = true
	log.Trace(errors.New("vPoint.Start()"))
	vPoint.Start()
	v.vpoint = vPoint

	log.Trace(errors.New("v.escortingUP()"))
	v.escortingUP()

	//Now, surpress interrupt signal for 5 sec

	v.interuptDeferto = 1

	go func() {
		time.Sleep(5 * time.Second)
		v.interuptDeferto = 0
	}()

	v.vpnSetup()

	if v.conf != nil {
		env := v.conf.additionalEnv
		log.Trace(errors.New("Exec Upscript()"))
		err = v.runbash(v.conf.upscript, env)
		if err != nil {
			log.Trace(errors.New("OnUp failed to exec").Base(err))
		}
	}

	v.Callbacks.OnEmitStatus(0, "Running")
	v.parseCfgDone()
}

/*RunLoop Run V2Ray main loop
 */
func (v *V2RayPoint) RunLoop() {
	v.v2rayOP.Lock()
	if !v.IsRunning {
		go v.pointloop()
	}
	v.v2rayOP.Unlock()
}

func (v *V2RayPoint) stopLoopW() {
	v.IsRunning = false
	v.vpoint.Close()

	if v.conf != nil {
		env := v.conf.additionalEnv
		log.Trace(errors.New("Running downscript"))
		err := v.runbash(v.conf.downscript, env)

		if err != nil {
			log.Trace(errors.New("OnDown failed to exec").Base(err))
		}
		log.Trace(errors.New("v.escortingDown()"))
		v.escortingDown()
	}

	v.Callbacks.OnEmitStatus(0, "Closed")

}

/*StopLoop Stop V2Ray main loop
 */
func (v *V2RayPoint) StopLoop() {
	v.v2rayOP.Lock()
	if v.IsRunning {
		v.vpnShutdown()
		go v.stopLoopW()
	}
	v.v2rayOP.Unlock()
}

/*NewV2RayPoint new V2RayPoint*/
func NewV2RayPoint() *V2RayPoint {
	return &V2RayPoint{unforgivnesschan: make(chan int), v2rayOP: new(sync.Mutex)}
}

/*NetworkInterrupted inform us to restart the v2ray,
closing dead connections.
*/
func (v *V2RayPoint) NetworkInterrupted() {
	/*
		Behavior Changed in API Ver 23
		From now, we will defer the start for 3 sec,
		any Interruption Message will be surpressed during this period
	*/
	go func() {
		if v.IsRunning {
			//Calc sleep time
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Your device might not support atomic operation", r)
				}
			}()
			log.Trace(errors.New("Running+NetworkInterrupted"))
			succ := atomic.CompareAndSwapInt64(&v.interuptDeferto, 0, 1)
			if succ {
				log.Trace(errors.New("Entered+NetworkInterrupted"))
				v.vpoint.Close()
				log.Trace(errors.New("Closed+NetworkInterrupted"))
				time.Sleep(2 * time.Second)
				log.Trace(errors.New("SleepDone+NetworkInterrupted"))
				v.vpoint.Start()
				log.Trace(errors.New("Started+NetworkInterrupted"))
				atomic.StoreInt64(&v.interuptDeferto, 0)
			} else {
				log.Trace(errors.New("Skipped+NetworkInterrupted"))
			}
		}
	}()
}

/*
Client can opt-in V2Ray's Next Generation Interface
*/
func (v *V2RayPoint) UpgradeToContext() {
	if v.Context == nil {
		v.Context = new(V2RayContext)
	}

}
