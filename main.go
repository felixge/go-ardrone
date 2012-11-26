package main

import (
	"github.com/felixge/ardrone/commands"
	"time"
	"fmt"
	"net"
)


func main() {
	sequence := new(commands.Sequence)
	conn, _ := net.Dial("udp", "192.168.1.1:5556")

	for {
		sequence.Add(&commands.Ref{Fly: true})
		sequence.Add(&commands.Pcmd{})

		message := sequence.ReadMessage()
		conn.Write([]byte(message))

		fmt.Printf("%#v\n", message)

		time.Sleep(30 * time.Millisecond)
	}
}
