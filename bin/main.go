package main

import (
	"fmt"
	"github.com/felixge/ardrone"
	/*"net"*/
	/*"time"*/
)

func main() {
	client, err := ardrone.Dial()

	if err != nil {
		panic(err)
	}

	fmt.Printf("client: %+v\n", client)

	/*sequence := new(commands.Sequence)*/
	/*controlConn, _ := net.Dial("udp", "192.168.1.1:5556")*/

	/*conn, err := navdata.Dial()*/
	/*if err != nil {*/
	/*panic(err)*/
	/*}*/

	/*for {*/
	/*data, err := conn.Read()*/
	/*if err != nil {*/
	/*fmt.Printf("error: %v\n", err)*/
	/*continue*/
	/*}*/

	/*fmt.Printf("%+v\n", data.Demo)*/

	/*sequence.Add(&commands.Ref{Emergency: true})*/
	/*[>sequence.Add(&commands.Pcmd{})<]*/

	/*message := sequence.ReadMessage()*/
	/*controlConn.Write([]byte(message))*/

	/*[>fmt.Printf("%#v\n", message)<]*/

	/*time.Sleep(30 * time.Millisecond)*/
	/*}*/
}
