package libv2ray

import (
	"os"
	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/app/router/rules"
	"v2ray.com/core/common/log"
	"v2ray.com/core/shell/point"

	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/proxy/blackhole"
	_ "v2ray.com/core/proxy/dokodemo"
	_ "v2ray.com/core/proxy/freedom"
	_ "v2ray.com/core/proxy/http"
	_ "v2ray.com/core/proxy/shadowsocks"
	_ "v2ray.com/core/proxy/socks"
	_ "v2ray.com/core/proxy/vmess/inbound"
	_ "v2ray.com/core/proxy/vmess/outbound"

	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/transport/internet/kcp"
	_ "v2ray.com/core/transport/internet/tcp"
	_ "v2ray.com/core/transport/internet/udp"
	_ "v2ray.com/core/transport/internet/ws"

	// The following are necessary as they register handlers in their init functions.
	_ "v2ray.com/core/transport/internet/authenticators/noop"
	_ "v2ray.com/core/transport/internet/authenticators/srtp"
	_ "v2ray.com/core/transport/internet/authenticators/utp"
)

/*V2RayPoint V2Ray Point Server
This is territory of Go, so no getter and setters!
*/
type V2RayPoint struct {
	ConfigureFile    string
	Callbacks        V2RayCallbacks
	vpoint           *point.Point
	IsRunning        bool
	conf             *libv2rayconf
	escortProcess    *[](*os.Process)
	unforgivnesschan chan int
	VpnSupportSet    V2RayVPNServiceSupportsSet
	VpnSupportnodup  bool
	PackageName      string
}

/*V2RayCallbacks a Callback set for V2Ray
 */
type V2RayCallbacks interface {
	OnEmitStatus(int, string) int
}

func (v *V2RayPoint) pointloop() {
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

	config, err := point.LoadConfig(v.ConfigureFile)
	if err != nil {
		log.Error("Failed to read config file (", v.ConfigureFile, "): ", v.ConfigureFile, err)

		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")

		return
	}
	if config.LogConfig != nil && len(config.LogConfig.AccessLog) > 0 {
		log.InitAccessLogger(config.LogConfig.AccessLog)
	}

	vPoint, err := point.NewPoint(config)
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
}

/*RunLoop Run V2Ray main loop
 */
func (v *V2RayPoint) RunLoop() {
	go v.pointloop()
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
	v.vpnShutdown()
	go v.stopLoopW()
}

/*NewV2RayPoint new V2RayPoint*/
func NewV2RayPoint() *V2RayPoint {
	return &V2RayPoint{unforgivnesschan: make(chan int)}
}

/*NetworkInterrupted inform us to restart the v2ray,
closing dead connections.
*/
func (v *V2RayPoint) NetworkInterrupted() {
	v.vpoint.Close()
	v.vpoint.Start()
}
