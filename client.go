package ardrone

import (
	"github.com/felixge/ardrone/navdata"
	"github.com/felixge/ardrone/commands"
	"fmt"
	"net"
	"time"
)

type Client struct {
	Config *Config
	navdataConn *navdata.Conn
	commands *commands.Sequence
	controlConn net.Conn
}

type Config struct {
	Ip string
	NavdataPort int
	AtPort int
	NavdataTimeout time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		Ip: "192.168.1.1",
		NavdataPort: 5554,
		AtPort: 5556,
		NavdataTimeout: 2000 * time.Millisecond,
	}
}

func (client *Client) Connect() error {
	navdataAddr := addr(client.Config.Ip, client.Config.NavdataPort)
	navdataConn, err := navdata.Dial(navdataAddr)
	if err != nil {
		return err
	}

	client.navdataConn = navdataConn
	client.navdataConn.SetReadTimeout(client.Config.NavdataTimeout)

	controlAddr := addr(client.Config.Ip, client.Config.AtPort)
	controlConn, err := net.Dial("udp", controlAddr)

	if err != nil {
		return err
	}

	client.controlConn = controlConn
	client.commands = &commands.Sequence{}


	ch := make(chan error)

	go func() { ch <- client.RequestDemoNavdata() }()
	go func() { ch <- client.DisableEmergency() }()


	for i := 0; i < 2; i++{
		if err = <-ch; err != nil {
			return err
		}
	}

	return nil
}

func (client *Client) RequestDemoNavdata() error{
	for {
		navdata, err := client.ReadNavdata()
		if err != nil {
			return err
		}

		if navdata.Demo != nil {
			return nil
		}

		client.commands.Add(commands.Config{Key: "general:navdata_demo", Value: "TRUE"})
		client.Send()
	}

	return nil
}

func (client *Client) DisableEmergency() error {
	for {
		data, err := client.ReadNavdata()
		if err != nil {
			return err
		}

		if data.Header.State & navdata.STATE_EMERGENCY_LANDING == 0 {
			return nil
		}

		client.commands.Add(&commands.Ref{Emergency: true})
		client.Send()
	}

	return nil
}

func (client *Client) ReadNavdata() (navdata *navdata.Navdata, err error) {
	return client.navdataConn.ReadNavdata()
}

func (client *Client) Takeoff() error {
	for {
		data, err := client.ReadNavdata()
		if err != nil {
			return err
		}

		if data.Demo.ControlState == navdata.CONTROL_HOVERING {
			break
		}

		client.commands.Add(&commands.Ref{Fly: true})
		client.commands.Add(&commands.Pcmd{})
		client.Send()
	}

	return nil
}

func (client *Client) Land() error {
	for {
		data, err := client.ReadNavdata()
		if err != nil {
			return err
		}

		if data.Demo.ControlState == navdata.CONTROL_LANDED {
			break
		}

		client.commands.Add(&commands.Ref{Fly: false})
		client.commands.Add(&commands.Pcmd{})
		client.Send()
	}

	return nil
}


func (client *Client) Send() {
	message := client.commands.ReadMessage()
	//fmt.Printf("message: %#v\n", message)
	client.controlConn.Write([]byte(message))
}

func addr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
