package ardrone

import (
	"github.com/felixge/ardrone/navdata"
)

func Dial() (client Client, err error) {
	client = Client{}

	navdataConn, err := navdata.Dial()
	if err != nil {
		return
	}

	client.navdataConn = navdataConn


	return
}

type Client struct{
	navdataConn navdata.Conn
}
