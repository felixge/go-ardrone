package main

import (
	"log"
	"github.com/felixge/ardrone"
	/*"net"*/
	"time"
)

func main() {
	log.SetFlags(log.Lmicroseconds)
	client := &ardrone.Client{Config: ardrone.DefaultConfig()}

	start := time.Now()

	log.Printf("Connecting to: %+v ...\n", client)

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Ready! Took %s\n", time.Since(start))

	start = time.Now()

	err = client.Takeoff()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Takeoff %s\n", time.Since(start))

	start = time.Now()

	err = client.Land()
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Land %s\n", time.Since(start))


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
