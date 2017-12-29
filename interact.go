package libv2ray

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"v2ray.com/core"
	"v2ray.com/ext/sysio"

	"github.com/golang/protobuf/proto"
	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process"
	"github.com/xiaokangwang/AndroidLibV2ray/Process/Escort"
	"github.com/xiaokangwang/AndroidLibV2ray/Process/UpDownScript"
	"github.com/xiaokangwang/AndroidLibV2ray/VPN"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
	"github.com/xiaokangwang/AndroidLibV2ray/configure/jsonConvert"
	"github.com/xiaokangwang/AndroidLibV2ray/shippedBinarys"
	vlencoding "github.com/xiaokangwang/V2RayConfigureFileUtil/encoding"
	mobasset "golang.org/x/mobile/asset"
	v2rayconf "v2ray.com/ext/tools/conf/serial"
)

/*V2RayPoint V2Ray Point Server
This is territory of Go, so no getter and setters!

Notice:
ConfigureFile can be either the path of config file or
"V2Ray_internal/ConfigureFileContent" in case you wish to
provide configure file with @[ConfigureFileContent] in JSON
format or "V2Ray_internal/AsPbConfigureFileContent"
in protobuf format.

*/
type V2RayPoint struct {
	status          *CoreI.Status
	confng          *configure.LibV2RayConf
	EnvCreater      Process.EnvironmentCreater
	escorter        *Escort.Escorting
	Callbacks       V2RayCallbacks
	v2rayOP         *sync.Mutex
	Context         *V2RayContext
	VPNSupports     *VPN.VPNSupport
	UpdownScripts   *UpDownScript.UpDownScript
	interuptDeferto int64

	//Legacy prop, should use Context instead
	PackageName          string
	ConfigureFile        string
	ConfigureFileContent string
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
	//Deal with legacy API
	var config core.Config
	if strings.HasPrefix(v.ConfigureFile, "V2Ray_internal") {
		if v.ConfigureFile == "V2Ray_internal/ConfigureFileContent" {
			//Convert is needed
			jc := &jsonConvert.JsonToPbConverter{}
			jc.Datadir = v.status.PackageName
			//Load File From Context
			//cf := v.Context.GetConfigureFile()
			jc.LoadFromString(v.ConfigureFileContent)
			jc.Parse()
			v.confng = jc.ToPb()
			configx, _ := v2rayconf.LoadJSONConfig(strings.NewReader(v.ConfigureFileContent))
			config = *configx
		} else {
			buf := []byte(v.ConfigureFileContent)
			err := proto.Unmarshal(buf, &config)
			if err != nil {
				log.Println(err)
			}
			//Assert V2RayPart
			for _, a := range config.GetExtension() {
				d, _ := a.GetInstance()
				switch vn := d.(type) {
				case *configure.LibV2RayConf:
					v.confng = vn
				default:
				}
			}
		}
	} else {
		//First Guess File type
		Type, err := vlencoding.GuessConfigType(v.Context.GetConfigureFile())
		if err != nil {
			fmt.Println(err)
			return
		}

		if Type == vlencoding.LibV2RayPackedConfig_FullJsonFile {
			//Convert is needed
			jc := &jsonConvert.JsonToPbConverter{}
			jc.Datadir = v.status.PackageName
			//Load File From Context
			cf := v.Context.GetConfigureFile()
			jc.LoadFromFile(cf)
			jc.Parse()
			v.confng = jc.ToPb()
			jsonctx, _ := os.Open(cf)
			configx, err := v2rayconf.LoadJSONConfig(jsonctx)
			if err != nil {
				log.Println("JSON Parse Err:" + err.Error())

			}
			if configx != nil {
				config = *configx
			}
			jsonctx.Close()
		} else if Type == vlencoding.LibV2RayPackedConfig_FullProto {
			buf, _ := ioutil.ReadFile(v.Context.GetConfigureFile())
			err = proto.Unmarshal(buf, &config)
			//Assert V2RayPart
			for _, a := range config.GetExtension() {
				d, _ := a.GetInstance()
				switch vn := d.(type) {
				case *configure.LibV2RayConf:
					v.confng = vn
				default:
				}
			}
		} else {

			//Yet To Support
			return
		}
	}
	var err error
	//TODO:Load Shipped Binary

	shipb := shippedBinarys.FirstRun{}
	shipb.SetCoreI(v.status)
	err = shipb.CheckAndExport()
	if err != nil {
		log.Println(err)
	}

	/*TODO:Load Client Config
	config, err := v2rayconf.LoadJSONConfig(v.parseCfg())
	if err != nil {
		log.Trace(errors.New("Failed to read config file (", v.ConfigureFile, "): ", v.ConfigureFile).Base(err).AtError())

		v.Callbacks.OnEmitStatus(-1, "Failed to read config file ("+v.ConfigureFile+")")

		return
	}*/

	v.status.Vpoint, err = core.New(&config)
	if err != nil {
		log.Println("VPoint Start Err:" + err.Error())

	}
	/* TODO: Start V2Ray Core
	vPoint, err := core.New(config)
	if err != nil {
		log.Trace(errors.New("Failed to create Point server").Base(err))

		v.Callbacks.OnEmitStatus(-1, "Failed to create Point server ("+err.Error()+")")

		return
	}*/

	v.status.IsRunning = true
	v.status.Vpoint.Start()
	/*log.Trace(errors.New("vPoint.Start()"))
	vPoint.Start()
	v.vpoint = vPoint
	*/
	if v.confng != nil {
		v.escorter.Configure = v.confng.RootModeConf.Escorting
		v.escorter.EscortingUP()

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
		//Set Necessary Props First
		v.VPNSupports.Conf = *v.confng.GetVpnConf()
		v.VPNSupports.SetStatus(v.status)
		v.VPNSupports.VpnSetup()
		/* TODO: setup VPN
		v.vpnSetup()
		*/
		v.UpdownScripts.SetStatus(v.status)
		v.UpdownScripts.Configure = v.confng.RootModeConf.Scripts
		v.UpdownScripts.Env = v.confng.Env
		v.UpdownScripts.RunUpScript()
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
	}
	v.Callbacks.OnEmitStatus(0, "Running")
	//v.parseCfgDone()
}

/*RunLoop Run V2Ray main loop
 */
func (v *V2RayPoint) RunLoop() {
	v.v2rayOP.Lock()
	//Construct Context
	if v.Context == nil {
		v.Context = new(V2RayContext)
		v.status.PackageName = v.PackageName
	}
	if !v.status.IsRunning {
		go v.pointloop()
	}
	v.v2rayOP.Unlock()
}

func (v *V2RayPoint) stopLoopW() {
	v.status.IsRunning = false
	v.status.Vpoint.Close()
	if v.confng != nil {
		v.UpdownScripts.RunDownScript()
		/* TODO:Run Down Script
		if v.conf != nil {
			env := v.conf.additionalEnv
			log.Trace(errors.New("Running downscript"))
			err := v.runbash(v.conf.downscript, env)

			if err != nil {
				log.Trace(errors.New("OnDown failed to exec").Base(err))
			}*/
		v.VPNSupports.VpnShutdown()
		v.escorter.EscortingDown()
		/* TODO: Escort Down
			log.Trace(errors.New("v.escortingDown()"))
			v.escortingDown()
		}
		*/
	}
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
	//Initialize asset API, Since Raymond Will not let notify the asset location inside Process,
	//We need to set location outside V2Ray
	const assetperfix = "/dev/libv2rayfs0/asset"
	os.Setenv("v2ray.location.asset", assetperfix)
	//Now we handle read
	sysio.NewFileReader = func(path string) (io.ReadCloser, error) {
		if strings.HasPrefix(path, assetperfix) {
			p := path[len(assetperfix)+1:]
			//is it overridden?
			by, ok := overridedAssets[p]
			if ok {
				return os.Open(by)
			}
			return mobasset.Open(p)
		}
		return os.Open(path)
	}
	//platform.ForceReevaluate()
	//panic("Creating VPoint")
	return &V2RayPoint{v2rayOP: new(sync.Mutex), status: &CoreI.Status{}, escorter: Escort.NewEscort(), VPNSupports: &VPN.VPNSupport{}, UpdownScripts: &UpDownScript.UpDownScript{}}
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
		v.Context.Status = v.status
	}
}

func (v *V2RayPoint) GetIsRunning() bool {
	return v.status.IsRunning
}

/*V2RayVPNServiceSupportsSet To support Android VPN mode*/
type V2RayVPNServiceSupportsSet interface {
	GetVPNFd() int
	Setup(Conf string) int
	Prepare() int
	Shutdown() int
	Protect(int) int
}

//Delegate Funcation
func (v *V2RayPoint) VpnSupportReady() {
	v.VPNSupports.VpnSupportReady()
}

//Delegate Funcation
func (v *V2RayPoint) SetVpnSupportSet(vs V2RayVPNServiceSupportsSet) {
	v.VPNSupports.VpnSupportSet = vs
}

func (v *V2RayPoint) OptinNextGenerationTunInterface() {
	v.VPNSupports.OptinNextGenerationTunInterface()
}
