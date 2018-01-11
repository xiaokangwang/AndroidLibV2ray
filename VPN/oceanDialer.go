package VPN

import (
	"context"
	"log"
	"net"
	"strconv"

	"v2ray.com/core"
	v2net "v2ray.com/core/common/net"
)

type V2Dialer struct {
	ser *core.Instance
}

func (vd *V2Dialer) Dial(network, address string, port uint16, ctx context.Context) (net.Conn, error) {
	var dest net.Addr
	var err error
	switch network {
	case "tcp4":
		dest, err = net.ResolveTCPAddr(network, address+":"+strconv.Itoa(int(port)))
		log.Println(err)
	case "udp4":
		dest, err = net.ResolveUDPAddr(network, address+":"+strconv.Itoa(int(port)))
		log.Println(err)
	}
	v2dest := v2net.DestinationFromAddr(dest)
	return core.Dial(ctx, vd.ser, v2dest)
}

func (vd *V2Dialer) NotifyMeltdown(reason error) {}
