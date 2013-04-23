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

	log.Printf("Connecting to: %+v ...\n", client)
	err := client.Connect()
	if err != nil {
		log.Fatal(err) // triggers a panic
	}

	log.Printf("Connected\n")
	client.Takeoff()
	log.Printf("Takeoff\n")
	client.ApplyFor(1500 * time.Millisecond, ardrone.State{Pitch: 0.5})
	time.Sleep(3 * time.Second)
	client.ApplyFor(2500 * time.Millisecond, ardrone.State{Pitch: -0.5})
	time.Sleep(3 * time.Second)
	client.Land()
	log.Printf("Land\n")
}
