package net

import (
	gonet "net"
	"testing"
	"time"
)

func TestNewPaVEDecoder(t *testing.T) {
	videoAddr := gonet.JoinHostPort(IP, VideoPort)
	conn, err := gonet.DialTimeout("tcp", videoAddr, time.Second*3)
	if err != nil {
		t.Skip(err)
	}

	f := PaVEFrame{}
	d := NewPaVEDecoder(conn)
	for i := 0; i < 3; i++ {
		err = d.Decode(&f)
		if err != nil {
			t.Fatal(err)
		}
		if string(f.Signature[0:]) != "PaVE" {
			t.Fatalf("Bad signature: %s", f.Signature)
		}
	}
}
