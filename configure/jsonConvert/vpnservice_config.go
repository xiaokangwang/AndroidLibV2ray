package jsonConvert

type vpnserviceConfig struct {
	Target      string   `json:"Target"`
	Args        []string `json:"Args"`
	VPNSetupArg string   `json:"VPNSetupArg"`
}

type vpnserviceDnsloopFix struct {
	DomainNameList []string `json:"domainName"`
	TCPVersion     string   `json:"tcpVersion"`
	UDPVersion     string   `json:"udpVersion"`
}
