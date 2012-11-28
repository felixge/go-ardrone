package navdata

import (
	"errors"
	"fmt"
	"net"
	"time"
)

var readTimeout = 1000 * time.Millisecond
var requestInterval = 100 * time.Millisecond

var ErrReadTimeout = errors.New("Conn.Read() timeout after " + fmt.Sprintf("%s", readTimeout))

type Conn struct {
	conn net.PacketConn
	addr *net.UDPAddr
}

func Dial() (conn Conn, err error) {
	udpConn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return
	}

	conn.conn = udpConn

	dst, err := net.ResolveUDPAddr("udp", "192.168.1.1:5554")
	if err != nil {
		return
	}

	conn.addr = dst
	return
}

func (this *Conn) Read() (navdata Data, err error) {
	timeout := time.After(readTimeout)
	write := time.NewTicker(requestInterval)
	defer write.Stop()

	read := make(chan []byte)

	go (func() {
		// The practical limit for the data length which is imposed by the
		// underlying IPv4 protocol is 65,507 bytes (65,535 − 8 byte UDP header −
		// 20 byte IP header).
		// -- http://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
		buf := make([]byte, 65507)
		n, _, _ := this.conn.ReadFrom(buf)
		read <- buf[:n]
	})()

	for {
		select {
		case data := <-read:
			navdata, err = Parse(data)
			return
		case <-write.C:
			output := []byte("\x01")
			_, err = this.conn.WriteTo(output, this.addr)
			if err != nil {
				return
			}
		case <-timeout:
			err = ErrReadTimeout
			return
		}
	}

	return
}
