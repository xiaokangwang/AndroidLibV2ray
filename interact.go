package libv2ray

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
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

	status          *CoreI.Status
	confng          *configure.LibV2RayConf
	EnvCreater      Process.EnvironmentCreater
	Callbacks       V2RayCallbacks
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

	v.status.VpnSupportnodup = false

	//TODO:Parse Configure File

	//TODO:Load Shipped Binary

	/*TODO:Load Client Config
	config, err := v2rayconf.LoadJSONConfig(v.parseCfg())
	if err != nil {
		log.Trace(errors.New("Failed to read config file (", v.ConfigureFile, "): ", v.ConfigureFile).Base(err).AtError())

		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")

		return
	}*/
	/* TODO: Start V2Ray Core
	vPoint, err := core.New(config)
	if err != nil {
		log.Trace(errors.New("Failed to create Point server").Base(err))

		v.Callbacks.OnEmitStatus(-1, "Failed to create Point server ("+err.Error()+")")

		return
	}*/

	v.status.IsRunning = true
	/*log.Trace(errors.New("vPoint.Start()"))
	vPoint.Start()
	v.vpoint = vPoint
	*/
	/* TODO:RunVPN Escort
	log.Trace(errors.New("v.escortingUP()"))
	v.escortingUP()
	*/
	//Now, surpress interrupt signal for 5 sec

	v.interuptDeferto = 1

	go func() {
		time.Sleep(5 * time.Second)
		v.interuptDeferto = 0
	}()
	/* TODO: setup VPN
	v.vpnSetup()
	*/
	/* TODO: Run Up Script
	if v.conf != nil {
		env := v.conf.additionalEnv
		log.Trace(errors.New("Exec Upscript()"))
		err = v.runbash(v.conf.upscript, env)
		if err != nil {
			log.Trace(errors.New("OnUp failed to exec").Base(err))
		}
	}
	*/
	v.Callbacks.OnEmitStatus(0, "Running")
	//v.parseCfgDone()
}

/*RunLoop Run V2Ray main loop
 */
func (v *V2RayPoint) RunLoop() {
	v.v2rayOP.Lock()
	if !v.status.IsRunning {
		go v.pointloop()
	}
	v.v2rayOP.Unlock()
}

func (v *V2RayPoint) stopLoopW() {
	v.status.IsRunning = false
	v.status.Vpoint.Close()
	/* TODO:Run Down Script
	if v.conf != nil {
		env := v.conf.additionalEnv
		log.Trace(errors.New("Running downscript"))
		err := v.runbash(v.conf.downscript, env)

		if err != nil {
			log.Trace(errors.New("OnDown failed to exec").Base(err))
		}*/
	/* TODO: Escort Down
		log.Trace(errors.New("v.escortingDown()"))
		v.escortingDown()
	}
	*/
	v.Callbacks.OnEmitStatus(0, "Closed")

}

/*StopLoop Stop V2Ray main loop
 */
func (v *V2RayPoint) StopLoop() {
	v.v2rayOP.Lock()
	if v.status.IsRunning {
		/* TODO: Shutdown VPN
		v.vpnShutdown()
		*/
		go v.stopLoopW()
	}
	v.v2rayOP.Unlock()
}

/*NewV2RayPoint new V2RayPoint*/
func NewV2RayPoint() *V2RayPoint {
	return &V2RayPoint{v2rayOP: new(sync.Mutex)}
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
		if v.status.IsRunning {
			//Calc sleep time
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Your device might not support atomic operation", r)
				}
			}()
			succ := atomic.CompareAndSwapInt64(&v.interuptDeferto, 0, 1)
			if succ {
				v.status.Vpoint.Close()
				time.Sleep(2 * time.Second)
				v.status.Vpoint.Start()
				atomic.StoreInt64(&v.interuptDeferto, 0)
			} else {
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

/*V2RayVPNServiceSupportsSet To support Android VPN mode*/
type V2RayVPNServiceSupportsSet interface {
	GetVPNFd() int
	Setup(Conf string) int
	Prepare() int
	Shutdown() int
	Protect(int) int
}
