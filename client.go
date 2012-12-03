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

func (this *Client) Connect() error {
	navdataAddr := addr(this.Config.Ip, this.Config.NavdataPort)
	navdataConn, err := navdata.Dial(navdataAddr)
	if err != nil {
		return err
	}

	this.navdataConn = navdataConn
	this.navdataConn.SetReadTimeout(this.Config.NavdataTimeout)

	controlAddr := addr(this.Config.Ip, this.Config.AtPort)
	controlConn, err := net.Dial("udp", controlAddr)

	if err != nil {
		return err
	}

	this.controlConn = controlConn
	this.commands = &commands.Sequence{}


	ch := make(chan error)

	go func() { ch <- this.RequestDemoNavdata() }()
	go func() { ch <- this.DisableEmergency() }()


	for i := 0; i < 2; i++{
		if err = <-ch; err != nil {
			return err
		}
	}

	return nil
}

func (this *Client) RequestDemoNavdata() error{
	for {
		navdata, err := this.ReadNavdata()
		if err != nil {
			return err
		}

		if navdata.Demo != nil {
			return nil
		}

		this.commands.Add(commands.Config{Key: "general:navdata_demo", Value: "TRUE"})
		this.Send()
	}

	return nil
}

func (this *Client) DisableEmergency() error {
	for {
		data, err := this.ReadNavdata()
		if err != nil {
			return err
		}

		if data.Header.State & navdata.STATE_EMERGENCY_LANDING == 0 {
			return nil
		}

		this.commands.Add(&commands.Ref{Emergency: true})
		this.Send()
	}

	return nil
}

func (this *Client) ReadNavdata() (navdata *navdata.Navdata, err error) {
	return this.navdataConn.ReadNavdata()
}

func (this *Client) Takeoff() error {
	for {
		data, err := this.ReadNavdata()
		if err != nil {
			return err
		}

		if data.Demo.ControlState == navdata.CONTROL_HOVERING {
			break
		}

		this.commands.Add(&commands.Ref{Fly: true})
		this.commands.Add(&commands.Pcmd{})
		this.Send()
	}

	return nil
}

func (this *Client) Land() error {
	for {
		data, err := this.ReadNavdata()
		if err != nil {
			return err
		}

		if data.Demo.ControlState == navdata.CONTROL_LANDED {
			break
		}

		this.commands.Add(&commands.Ref{Fly: false})
		this.commands.Add(&commands.Pcmd{})
		this.Send()
	}

	return nil
}


func (this *Client) Send() {
	message := this.commands.ReadMessage()
	//fmt.Printf("message: %#v\n", message)
	this.controlConn.Write([]byte(message))
}

func addr(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)
}
