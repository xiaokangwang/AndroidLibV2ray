package VPN

import (
	"log"
	"net"

	"github.com/davecgh/go-spew/spew"
)

type preparedDomain struct {
	tcpprepared map[string](*net.TCPAddr)
	udpprepared map[string](*net.UDPAddr)
}

func (v *VPNSupport) prepareDomainName() {
	if v.VpnSupportSet == nil {
		return
	}
	v.prepareddomain.tcpprepared = make(map[string](*net.TCPAddr))
	v.prepareddomain.udpprepared = make(map[string](*net.UDPAddr))
	for _, domainName := range v.Conf.PreparedDomainName.DomainNameList {
		log.Println("Preparing DNS,", domainName)
		var err error
		v.prepareddomain.tcpprepared[domainName], err = net.ResolveTCPAddr(v.Conf.PreparedDomainName.TCPVersion, domainName)
		if err != nil {
			log.Println(err)
		}
		v.prepareddomain.udpprepared[domainName], err = net.ResolveUDPAddr(v.Conf.PreparedDomainName.UDPVersion, domainName)
		spew.Dump(v.prepareddomain.udpprepared[domainName])
		if err != nil {
			log.Println(err)
		}
	}
}
