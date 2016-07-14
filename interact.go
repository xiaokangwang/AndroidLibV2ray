package libv2ray

import (
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
	ConfigureFile string
	OnEmitStatus  *func(StatusCode int, Status string) int
	vpoint        *point.Point
}

func (v *V2RayPoint) pointloop() {
	config, err := point.LoadConfig(v.ConfigureFile)
	if err != nil {
		log.Error("Failed to read config file (", v.ConfigureFile, "): ", v.ConfigureFile, err)
		if v.OnEmitStatus != nil {
			(*v.OnEmitStatus)(-1, "Failed to read config file ("+v.ConfigureFile+")")
		}
		return
	}
	if config.LogConfig != nil && len(config.LogConfig.AccessLog) > 0 {
		log.InitAccessLogger(config.LogConfig.AccessLog)
	}

	vPoint, err := point.NewPoint(config)
	if err != nil {
		log.Error("Failed to create Point server: ", err)
		if v.OnEmitStatus != nil {
			(*v.OnEmitStatus)(-1, "Failed to create Point server ("+err.Error()+")")
		}
		return
	}
	vPoint.Start()
	v.vpoint = vPoint
}

/*RunLoop Run V2Ray main loop
 */
func (v *V2RayPoint) RunLoop() {
	go v.pointloop()
}

/*StopLoop Stop V2Ray main loop
 */
func (v *V2RayPoint) StopLoop() {
	v.vpoint.Close()
}
