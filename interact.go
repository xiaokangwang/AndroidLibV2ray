package libv2ray

import (
	"os"
	// The following are necessary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/app/router/rules"
	"github.com/v2ray/v2ray-core/common/log"
	"github.com/v2ray/v2ray-core/shell/point"

	// The following are necessary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/proxy/blackhole"
	_ "github.com/v2ray/v2ray-core/proxy/dokodemo"
	_ "github.com/v2ray/v2ray-core/proxy/freedom"
	_ "github.com/v2ray/v2ray-core/proxy/http"
	_ "github.com/v2ray/v2ray-core/proxy/shadowsocks"
	_ "github.com/v2ray/v2ray-core/proxy/socks"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/inbound"
	_ "github.com/v2ray/v2ray-core/proxy/vmess/outbound"

	// The following are necessary as they register handlers in their init functions.
	_ "github.com/v2ray/v2ray-core/transport/internet/kcp"
	_ "github.com/v2ray/v2ray-core/transport/internet/tcp"
	_ "github.com/v2ray/v2ray-core/transport/internet/udp"
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
}

/*V2RayCallbacks a Callback set for V2Ray
 */
type V2RayCallbacks interface {
	OnEmitStatus(int, string) int
}

func (v *V2RayPoint) pointloop() {
	if v.parseConf() != nil {
		return
	}

	err := checkIfRcExist()

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
	go v.stopLoopW()
}

/*NewV2RayPoint new V2RayPoint*/
func NewV2RayPoint() *V2RayPoint {
	return &V2RayPoint{unforgivnesschan: make(chan int)}
}
