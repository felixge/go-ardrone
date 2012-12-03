package navdata

// @TODO: make this work again

import (
	"time"
//"errors"
//"fmt"
	"net"
)

type ErrReadTimeout time.Duration

func (this ErrReadTimeout) Error() string {
	// @TODO: Why does this not work?
	//return "navdata: read timeout after " + this.String()
	return "navdata: read timeout after "
}

type readCall interface{}
type readResult struct{
	navdata *Navdata
	err error
}

type Conn struct {
	udpConn net.PacketConn
	addr *net.UDPAddr
	readTimeout time.Duration
	readCall chan readCall
	readResult chan readResult
}

func Dial(addr string) (conn *Conn, err error) {
	conn = &Conn{
		readCall: make(chan readCall),
		readResult: make(chan readResult),
	}

	udpConn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return
	}

	conn.udpConn = udpConn

	dst, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return
	}

	conn.addr = dst
	go conn.readFan()

	return
}

func (this *Conn) SetReadTimeout(timeout time.Duration) {
	this.readTimeout = timeout
}

func (this *Conn) readFan() {
	for {
		// Sleep until we have somebody calling us
		<-this.readCall
		numCalls := 1

		// Ok, time to some work
		navdata, err := this.readNavdata()
		result := readResult{navdata, err}

	CollectCallers:
		for {
			select {
			case <-this.readCall:
				numCalls++
			default:
				break CollectCallers
			}
		}

		for i := 0; i < numCalls; i++ {
			this.readResult <- result
		}
	}
}


func (this *Conn) ReadNavdata() (navdata *Navdata, err error) {
	this.readCall <- nil
	result := <-this.readResult
	return result.navdata, result.err
}

var requestInterval = 100 * time.Millisecond

func (this *Conn) readNavdata() (navdata *Navdata, err error) {
	readTimeout := time.After(this.readTimeout)
	write := time.NewTicker(requestInterval)

	defer write.Stop()

	readResult := make(chan []byte)
	readError := make(chan error)

	go (func() {
		// The practical limit for the data length which is imposed by the
		// underlying IPv4 protocol is 65,507 bytes (65,535 − 8 byte UDP header −
		// 20 byte IP header).
		// -- http://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
		buf := make([]byte, 65507)
		n, _, err := this.udpConn.ReadFrom(buf)

		if err != nil {
			readError <- err
		} else {
			readResult <- buf[:n]
		}
	})()

	for {
		select {
		case data := <-readResult:
			navdata, err = Decode(data)
			return
		case <-write.C:
			_, err = this.udpConn.WriteTo([]byte("\x01"), this.addr)
			if err != nil {
				return
			}
		case <-readTimeout:
			err = ErrReadTimeout(this.readTimeout)
			return
		}
	}

	return
}
