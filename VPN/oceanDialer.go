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
	"v2ray.com/core/common/signal"
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

	input := ray.InboundInput()
	output := ray.InboundOutput()

	requestDone := signal.ExecuteAsync(func() error {
		defer input.Close()
		v2reader := buf.NewReader(conn1)
		if err := buf.Copy(v2reader, input); err != nil {
			return err
		}
		return nil
	})

	responseDone := signal.ExecuteAsync(func() error {
		v2writer := buf.NewWriter(conn1)
		if err := buf.Copy(output, v2writer); err != nil {
			return err
		}
		return nil
	})

	go func() {
		if err := signal.ErrorOrFinish2(ctx, requestDone, responseDone); err != nil {
			input.CloseError()
			output.CloseError()
			return
		}
	}()

	return conn2, nil
}

func (vd *V2Dialer) NotifyMeltdown(reason error) {}
