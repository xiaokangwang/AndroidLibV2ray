package VPN

import (
	"context"
	"os"
	"reflect"
	"syscall"
	"unsafe"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process/Escort"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"
	"github.com/xiaokangwang/waVingOcean"
	voconfigure "github.com/xiaokangwang/waVingOcean/configure"

	"golang.org/x/sys/unix"

	"v2ray.com/core/app"
	"v2ray.com/core/transport/internet"
)

/*VpnSupportReady VpnSupportReady*/
func (v *VPNSupport) VpnSupportReady() {
	if !v.status.VpnSupportnodup {
		/*
			v.VpnSupportnodup = true
			//Surpress Network Interruption Notifiction
			go func() {
				time.Sleep(5 * time.Second)
				v.VpnSupportnodup = false
			}()*/
		v.VpnSupportSet.Setup(v.Conf.Service.VPNSetupArg)
		v.setV2RayDialer()
		v.startVPNRequire()
	}
}
func (v *VPNSupport) startVPNRequire() {
	if !v.usewaVingOceanVPNBackend {
		v.Estr = Escort.NewEscort()
		v.Estr.SetStatus(v.status)
		v.Estr.EscortingUPV()
		go v.Estr.EscortRun(v.Conf.Service.Target, v.Conf.Service.Args, false, v.VpnSupportSet.GetVPNFd())
	} else {
		v.startNextGen()
	}
}

func (v *VPNSupport) askSupportSetInit() {
	v.VpnSupportSet.Prepare()
}

func (v *VPNSupport) VpnSetup() {
	if v.Conf.Service.VPNSetupArg != "" {
		v.prepareDomainName()

		v.askSupportSetInit()
	}
}
func (v *VPNSupport) VpnShutdown() {

	if v.Conf.Service.VPNSetupArg != "" {
		/*
			BUG DISCOVERED!

			v.VpnSupportnodup can have unexpected value cause VPN failed to revoke.
			more testing needed.

		*/

		//if v.VpnSupportnodup {
		err := unix.Close(v.VpnSupportSet.GetVPNFd())
		println(err)
		//}
		v.VpnSupportSet.Shutdown()
		if !v.usewaVingOceanVPNBackend {
			v.Estr.EscortingDown()
		} else {
			v.stopNextGen()
		}

	}
	v.status.VpnSupportnodup = false
}

func (v *VPNSupport) setV2RayDialer() {
	protectedDialer := &vpnProtectedDialer{vp: v}
	internet.UseAlternativeSystemDialer(internet.WithAdapter(protectedDialer))
}

type VPNSupport struct {
	prepareddomain           preparedDomain
	VpnSupportSet            V2RayVPNServiceSupportsSet
	status                   *CoreI.Status
	Conf                     configure.VPNConfig
	Estr                     *Escort.Escorting
	usewaVingOceanVPNBackend bool
	lowerup                  *wavingocean.LowerUp
}

func (v *VPNSupport) startNextGen() {
	tapfd := v.VpnSupportSet.GetVPNFd()
	syscall.SetNonblock(tapfd, false)
	f := os.NewFile(uintptr(tapfd), "/dev/tap0")
	cfg := new(voconfigure.WaVingOceanConfigure)
	cfg.PublicOnly = false
	cfg.EnableDnsCache = false
	cfg.DNSServers = make([]string, 0)
	v.lowerup = wavingocean.NewLowerUp(*cfg, f, &V2Dialer{ser: v.getSpace()}, context.TODO())
	go v.lowerup.Up()
}

func (v *VPNSupport) OptinNextGenerationTunInterface() {
	v.usewaVingOceanVPNBackend = true
}

func (v *VPNSupport) stopNextGen() {
	v.lowerup.Down()
}

func (v *VPNSupport) getSpace() app.Space {
	VpV := reflect.ValueOf(v.status.Vpoint)
	Space := VpV.Elem().FieldByName("space")
	//unsafely neutralize unexport field protection
	Space = reflect.NewAt(Space.Type(), unsafe.Pointer(Space.UnsafeAddr()))
	s := Space.Elem().Interface().(app.Space)
	return s
}

type V2RayVPNServiceSupportsSet interface {
	GetVPNFd() int
	Setup(Conf string) int
	Prepare() int
	Shutdown() int
	Protect(int) int
}

func (v *VPNSupport) SetStatus(st *CoreI.Status) {
	v.status = st
}
