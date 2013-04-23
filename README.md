# ardrone

A Go implementation of the Parrot AR Drone protocols. Not complete yet, but
the stuff that is implemented should work : ).


Get the latest version from Github
```bash
go get github.com/felixge/ardrone
```


The code below will execute a nice little sequence, but please make sure you
have enough space when running it.

```js
package main

import (
	"github.com/felixge/ardrone"
	"time"
)

func main() {
	client, err := ardrone.Connect(ardrone.DefaultConfig())
	if err != nil {
		panic(err)
	}

	client.Takeoff()
	client.Vertical(1*time.Second, 0.5)
	time.Sleep(3 * time.Second)
	client.Roll(1*time.Second, 0.5)
	time.Sleep(3 * time.Second)
	client.Roll(1*time.Second, -0.5)
	time.Sleep(3 * time.Second)
	client.Animate(ardrone.FLIP_LEFT, 200)
	time.Sleep(5 * time.Second)
	client.Land()
}
```
Save the code into a file and run:

```bash
go run main.go
```

## API Documentation

The API documentation can be found [here](http://godoc.org/github.com/felixge/ardrone).
