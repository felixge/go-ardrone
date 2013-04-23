# ardrone

A Go implementation of the Parrot AR Drone protocols. Not complete yet, but
the stuff that is implemented should work : ).


Get the latest version from Github
```bash
go get github.com/felixge/ardrone
```


Simple testcode to get the drone to takeoff, fly forward and land:

```js
package main

import (
  "github.com/felixge/ardrone"
  "time"
)

func main() {
  client, err := ardrone.Connect{ardrone.DefaultConfig()}
  if (err != nil) {
    panic(err)
  }

  client.Takeoff()
  client.ApplyFor(1 * time.Second, ardrone.State{Pitch: 0.5})
  time.Sleep(3 * time.Second)
  client.Land()
}
```
Save the code into a file and run:

```bash
go run main.go
```

## API Documentation

The API documentation can be found [here](http://godoc.org/github.com/felixge/ardrone).
