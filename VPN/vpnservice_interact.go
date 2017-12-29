package VPN

import (
	"reflect"

	"github.com/xiaokangwang/AndroidLibV2ray/CoreI"
	"github.com/xiaokangwang/AndroidLibV2ray/Process/Escort"
	"github.com/xiaokangwang/AndroidLibV2ray/configure"

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
	v.Estr = Escort.NewEscort()
	v.Estr.SetStatus(v.status)
	v.Estr.EscortingUPV()
	go v.Estr.EscortRun(v.Conf.Service.Target, v.Conf.Service.Args, false, v.VpnSupportSet.GetVPNFd())
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
		v.Estr.EscortingDown()
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
}

func (v *VPNSupport) getSpace() app.Space {
	VpV := reflect.ValueOf(v.status.Vpoint)
	Space := VpV.FieldByName("space")
	s := Space.Interface().(app.Space)
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
