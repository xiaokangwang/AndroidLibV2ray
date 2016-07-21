package libv2ray

import (
	"net"
	"time"
)

type surrogateDialer struct {
	Timeout       time.Duration
	Deadline      time.Time
	LocalAddr     net.Addr
	DualStack     bool
	FallbackDelay time.Duration
	KeepAlive     time.Duration
	Cancel        <-chan struct{}
}

/*
type surrogateConn struct {
}

func (sc *surrogateConn) Read(b []byte) (n int, err error) {
	return 0, nil
}
func (sc *surrogateConn) Write(b []byte) (n int, err error) {
	return 0, nil
}
func (sc *surrogateConn) Close() error {
	return nil
}
func (sc *surrogateConn) LocalAddr() net.Addr {
	return nil
}
func (sc *surrogateConn) RemoteAddr() net.Addr {
	return nil
}
func (sc *surrogateConn) SetDeadline(t time.Time) error {
	return nil
}
func (sc *surrogateConn) SetReadDeadline(t time.Time) error {
	return nil
}
func (sc *surrogateConn) SetWriteDeadline(t time.Time) error {
	return nil

}
*/

/*V2RayVPNServiceSupportsSet To support Android VPN mode*/
type V2RayVPNServiceSupportsSet interface {
	GetVPNFd() int
	Setup(Conf string) int
	Prepare() int
	Shutdown() int
}
