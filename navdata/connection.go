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

type Conn struct {
	udpConn net.PacketConn
	addr *net.UDPAddr
}

func Dial(addr string) (conn *Conn, err error) {
	conn = new(Conn)

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

	return
}

var requestInterval = 100 * time.Millisecond

func (this *Conn) ReadNavdata(timeout time.Duration) (navdata *Navdata, err error) {
	readTimeout := time.After(timeout)
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
			err = ErrReadTimeout(timeout)
			return
		}
	}

	return
}

//func (this *Conn) Read() (navdata Data, err error) {
//timeout := time.After(readTimeout)
//write := time.NewTicker(requestInterval)
//defer write.Stop()

//read := make(chan []byte)

//go (func() {
//// The practical limit for the data length which is imposed by the
//// underlying IPv4 protocol is 65,507 bytes (65,535 − 8 byte UDP header −
//// 20 byte IP header).
//// -- http://en.wikipedia.org/wiki/User_Datagram_Protocol#Packet_structure
//buf := make([]byte, 65507)
//n, _, _ := this.conn.ReadFrom(buf)
//read <- buf[:n]
//})()

//for {
//select {
//case data := <-read:
//navdata, err = Parse(data)
//return
//case <-write.C:
//output := []byte("\x01")
//_, err = this.conn.WriteTo(output, this.addr)
//if err != nil {
//return
//}
//case <-timeout:
//err = ErrReadTimeout
//return
//}
//}

//return
//}
