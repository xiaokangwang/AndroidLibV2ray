package libv2ray

/*VpnSupportReady VpnSupportReady*/
func (v *V2RayPoint) VpnSupportReady() {
	v.VpnSupportSet.Setup(v.conf.vpnConfig.VPNSetupArg)

}
func (v *V2RayPoint) startVPNRequire() {
	v.escortRun(v.conf.vpnConfig.Target, v.conf.vpnConfig.Args, false, v.VpnSupportSet.GetVPNFd())
}
func (v *V2RayPoint) askSupportSetInit() {
	v.VpnSupportSet.Prepare()
}
