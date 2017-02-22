package libv2ray

import (
	"os"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core"
	"v2ray.com/core/app/log"

	// For json config parser
	_ "v2ray.com/core/tools/conf"
)

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
	conf                 *libv2rayconf
	escortProcess        *[](*os.Process)
	unforgivnesschan     chan int
	VpnSupportSet        V2RayVPNServiceSupportsSet
	VpnSupportnodup      bool
	PackageName          string
	cfgtmpvarsi          cfgtmpvars
	//softcrashMonitor     bool
	prepareddomain  preparedDomain
	v2rayOP         *sync.Mutex
	interuptDeferto int64
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
		log.Error("Failed to copy asset", err)
		v.Callbacks.OnEmitStatus(-1, "Failed to copy asset ("+err.Error()+")")

	}

	log.Info("v.renderAll() ")
	v.renderAll()

	//Surpress Network Interruption Notifiction
	atomic.StoreInt64(&v.interuptDeferto, 1)

	config, err := core.LoadConfig(core.ConfigFormat_JSON, v.parseCfg())
	if err != nil {
		log.Error("Failed to read config file (", v.ConfigureFile, "): ", v.ConfigureFile, err)

		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")

		return
	}

	vPoint, err := core.New(config)
	if err != nil {
		log.Error("Failed to create Point server: ", err)

		v.Callbacks.OnEmitStatus(-1, "Failed to create Point server ("+err.Error()+")")

		return
	}
	v.IsRunning = true
	log.Info("vPoint.Start() ")
	vPoint.Start()
	v.vpoint = vPoint

	log.Info("v.escortingUP() ")
	v.escortingUP()

	v.vpnSetup()

	if v.conf != nil {
		env := v.conf.additionalEnv
		log.Info("Exec Upscript() ")
		err = v.runbash(v.conf.upscript, env)
		if err != nil {
			log.Error("OnUp failed to exec: ", err)
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
		log.Info("Running downscript")
		err := v.runbash(v.conf.downscript, env)

		if err != nil {
			log.Error("OnDown failed to exec: ", err)
		}
		log.Info("v.escortingDown() ")
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
			log.Info("Running+NetworkInterrupted")
			succ := atomic.CompareAndSwapInt64(&v.interuptDeferto, 0, 1)
			if succ {
				log.Info("Entered+NetworkInterrupted")
				v.vpoint.Close()
				log.Info("Closed+NetworkInterrupted")
				time.Sleep(3 * time.Second)
				log.Info("SleepDone+NetworkInterrupted")
				v.vpoint.Start()
				log.Info("Started+NetworkInterrupted")
				atomic.StoreInt64(&v.interuptDeferto, 0)
			} else {
				log.Info("Skipped+NetworkInterrupted")
			}
		}
	}()
}
