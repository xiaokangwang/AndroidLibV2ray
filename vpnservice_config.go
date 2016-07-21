package libv2ray

type vpnserviceConfig struct {
	Target      string   `json:"Target"`
	Args        []string `json:"Args"`
	VPNSetupArg string   `json:"VPNSetupArg"`
}
