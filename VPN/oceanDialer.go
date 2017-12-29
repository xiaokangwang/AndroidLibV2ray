package VPN

import (
	"context"
	"log"
	"net"
	"strconv"

	"v2ray.com/core/app"
	"v2ray.com/core/app/dispatcher"
	"v2ray.com/core/common/buf"
	v2net "v2ray.com/core/common/net"
)

type V2Dialer struct {
	ser app.Space
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
	disp := dispatcher.FromSpace(vd.ser)
	ray, err := disp.Dispatch(ctx, v2dest)
	if err != nil {
		panic(err)
	}
	//Copy data
	conn1, conn2 := net.Pipe()
	go func() {
		for {
			var buffer [1500]byte
			buf, err := ray.InboundOutput().ReadMultiBuffer()
			if err != nil {
				log.Println(err)
				return
			}
			n, err := buf.Read(buffer[:])
			if err != nil {
				log.Println(err)
				return
			}
			conn1.Write(buffer[:n])
		}
	}()
	go func() {
		for {
			mb := buf.NewMultiBufferCap(65536)
			var buffer [1500]byte
			n, err := conn1.Read(buffer[:])
			if err != nil {
				log.Println(err)
				return
			}
			mb.Write(buffer[:n])
			err = ray.InboundInput().WriteMultiBuffer(mb)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}()
	return conn2, nil
}

func (vd *V2Dialer) NotifyMeltdown(reason error) {}
