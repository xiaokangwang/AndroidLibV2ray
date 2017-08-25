package jsonConvert

import (
	"fmt"
	"strings"

	"github.com/xiaokangwang/AndroidLibV2ray/configure"
)

func ConvertToPb(leagcy libv2rayconf) *configure.LibV2RayConf {
	NextGenerationProtobufConfigureStruct := &configure.LibV2RayConf{}
	NextGenerationProtobufConfigureStruct.NoAutoConvert = true
	NextGenerationProtobufConfigureStruct.RootModeConf = &configure.RootModeConfig{}
	NextGenerationProtobufConfigureStruct.RootModeConf.Scripts = &configure.UpDownScripts{}
	NextGenerationProtobufConfigureStruct.RootModeConf.Scripts.UpScript = leagcy.upscript
	NextGenerationProtobufConfigureStruct.RootModeConf.Scripts.DownScript = leagcy.downscript
	NextGenerationProtobufConfigureStruct.RootModeConf.Escorting = make([]*configure.EscortedProcess, 0, len(leagcy.esco))
	for _, EscortedProcessInLegacy := range leagcy.esco {
		designatedAppendee := new(configure.EscortedProcess)
		designatedAppendee.Target = EscortedProcessInLegacy.Target
		designatedAppendee.Forgiveable = EscortedProcessInLegacy.Forgiveable
		designatedAppendee.Args = EscortedProcessInLegacy.Args
		NextGenerationProtobufConfigureStruct.RootModeConf.Escorting = append(NextGenerationProtobufConfigureStruct.RootModeConf.Escorting, designatedAppendee)
	}
	NextGenerationProtobufConfigureStruct.VpnConf = &configure.VPNConfig{}
	NextGenerationProtobufConfigureStruct.VpnConf.Service = &configure.VPNServiceConfig{}
	NextGenerationProtobufConfigureStruct.VpnConf.Service.Target = leagcy.vpnConfig.Target
	NextGenerationProtobufConfigureStruct.VpnConf.Service.VPNSetupArg = leagcy.vpnConfig.VPNSetupArg
	NextGenerationProtobufConfigureStruct.VpnConf.Service.Args = leagcy.vpnConfig.Args
	NextGenerationProtobufConfigureStruct.VpnConf.PreparedDomainName = &configure.DNSLoopFix{}
	NextGenerationProtobufConfigureStruct.VpnConf.PreparedDomainName.TCPVersion = leagcy.dnsloopfix.TCPVersion
	NextGenerationProtobufConfigureStruct.VpnConf.PreparedDomainName.UDPVersion = leagcy.dnsloopfix.UDPVersion
	NextGenerationProtobufConfigureStruct.VpnConf.PreparedDomainName.DomainNameList = leagcy.dnsloopfix.DomainNameList
	NextGenerationProtobufConfigureStruct.Env = &configure.EnvironmentVar{}
	NextGenerationProtobufConfigureStruct.Env.Vars = envToMap(leagcy.additionalEnv)
	fmt.Println(NextGenerationProtobufConfigureStruct.String())
	return NextGenerationProtobufConfigureStruct
}

func envToMap(k []string) map[string]string {
	var themap map[string]string
	themap = make(map[string]string)
	for _, val := range k {
		r := strings.Index(val, "=")
		themap[val[:r]] = val[r+1:]
	}
	return themap
}
