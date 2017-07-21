package libv2ray

import "os"

func (v *V2RayPoint) isDebugTriggered() bool {
	if _, err := os.Stat(v.Context.getDataDir() + "debug_enabled"); os.IsNotExist(err) {
		return false
	}
	return true
}

func (v *V2RayPoint) IsDebugTriggered() bool {
	return v.isDebugTriggered()
}
