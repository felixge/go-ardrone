ardrone
=======

**Warning:** This package is quite incomplete. If something is not working, it's probably not done yet. I do plan
to finish this up so!

A Go implementation of the Parrot AR Drone protocols.


Get the latest version from Github
```bash
go get github.com/felixge/ardrone
```


Simple testcode to get the drone to takeof and land:
```js
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
}
```
Save the code into a file and run 
```bash
go build main.go
```

Then run
```bash
./go
```


