package libv2ray

import "log"

/*VpnSupportReady VpnSupportReady*/
func (v *V2RayPoint) VpnSupportReady() {
	v.VpnSupportSet.Setup(v.conf.vpnConfig.VPNSetupArg)
	v.startVPNRequire()
}
func (v *V2RayPoint) startVPNRequire() {
	go v.escortRun(v.conf.vpnConfig.Target, v.conf.vpnConfig.Args, false, v.VpnSupportSet.GetVPNFd())
}
func (v *V2RayPoint) askSupportSetInit() {
	v.VpnSupportSet.Prepare()
}

func (v *V2RayPoint) vpnSetup() {
	log.Println(v.conf.vpnConfig.VPNSetupArg)
	if v.conf.vpnConfig.VPNSetupArg != "" {
		v.askSupportSetInit()
	}
}
func (v *V2RayPoint) vpnShutdown() {
	if v.conf.vpnConfig.VPNSetupArg != "" {
		v.VpnSupportSet.Shutdown()
	}
}
