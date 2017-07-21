package libv2ray

import (
	"log"
	"net"

	"github.com/davecgh/go-spew/spew"
)

type preparedDomain struct {
	tcpprepared map[string](*net.TCPAddr)
	udpprepared map[string](*net.UDPAddr)
}

func (v *V2RayPoint) prepareDomainName() {
	if v.conf == nil {
		return
	}
	v.prepareddomain.tcpprepared = make(map[string](*net.TCPAddr))
	v.prepareddomain.udpprepared = make(map[string](*net.UDPAddr))
	for _, domainName := range v.conf.dnsloopfix.DomainNameList {
		log.Println("Preparing DNS,", domainName)
		var err error
		v.prepareddomain.tcpprepared[domainName], err = net.ResolveTCPAddr(v.conf.dnsloopfix.TCPVersion, domainName)
		if err != nil {
			log.Println(err)
		}
		v.prepareddomain.udpprepared[domainName], err = net.ResolveUDPAddr(v.conf.dnsloopfix.UDPVersion, domainName)
		spew.Dump(v.prepareddomain.udpprepared[domainName])
		if err != nil {
			log.Println(err)
		}
	}
}
