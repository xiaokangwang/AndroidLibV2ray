package libv2ray

func (v *V2RayPoint) GetStatControler() *StatControler {
	return &StatControler{InterfaceTarget: "tun0:"}
}
